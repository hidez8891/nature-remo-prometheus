package natureremo

import (
	"context"
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
)

type Collector struct {
	natureToken    string
	natureUrl      string
	deviceMetrics  *prometheus.Desc
	echonetMetrics *prometheus.Desc
}

func NewCollector(natureToken, natureUrl string) *Collector {
	return &Collector{
		natureToken: natureToken,
		natureUrl:   natureUrl,
		deviceMetrics: prometheus.NewDesc(
			"natureremo_device_events",
			"NatureRemo Device Events",
			[]string{"device_name", "event"},
			nil,
		),
		echonetMetrics: prometheus.NewDesc(
			"natureremo_echonet_lite",
			"NatureRemo Echonet Lite",
			[]string{"device_name", "type"},
			nil,
		),
	}
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.deviceMetrics
	ch <- c.echonetMetrics
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	client, err := NewClient(c.natureUrl)
	if err != nil {
		ch <- prometheus.NewInvalidMetric(c.deviceMetrics, err)
		ch <- prometheus.NewInvalidMetric(c.echonetMetrics, err)
		return
	}

	// device events
	{
		ctx := context.Background()
		rsp, err := client.Get1Devices(ctx, c.addAuthHeader)
		if err != nil {
			ch <- prometheus.NewInvalidMetric(c.deviceMetrics, err)
			return
		}

		devices, err := ParseGet1DevicesResponse(rsp)
		if err != nil {
			ch <- prometheus.NewInvalidMetric(c.deviceMetrics, err)
			return
		}
		if devices.StatusCode() != 200 {
			err := fmt.Errorf("NatureRemo API return code=%d", devices.StatusCode())
			ch <- prometheus.NewInvalidMetric(c.deviceMetrics, err)
			return
		}

		for _, device := range *devices.JSON200 {
			for typ, event := range *device.NewestEvents {
				ch <- prometheus.NewMetricWithTimestamp(
					*event.CreatedAt,
					prometheus.MustNewConstMetric(
						c.deviceMetrics,
						prometheus.GaugeValue,
						float64(*event.Val),
						*device.Name,
						typ,
					),
				)
			}
		}
	}

	// Echonet Lite
	{
		ctx := context.Background()
		rsp, err := client.Get1EchonetliteAppliances(ctx, c.addAuthHeader)
		if err != nil {
			ch <- prometheus.NewInvalidMetric(c.echonetMetrics, err)
			return
		}

		appliances, err := ParseGet1EchonetliteAppliancesResponse(rsp)
		if err != nil {
			ch <- prometheus.NewInvalidMetric(c.echonetMetrics, err)
			return
		}
		if appliances.StatusCode() != 200 {
			err := fmt.Errorf("NatureRemo API return code=%d", appliances.StatusCode())
			ch <- prometheus.NewInvalidMetric(c.echonetMetrics, err)
			return
		}

		for _, appliance := range *appliances.JSON200.Appliances {
			echonet := NewEchonet()
			for _, prop := range *appliance.Properties {
				err := echonet.SetValue(*prop.Epc, *prop.Val, *prop.UpdatedAt)
				if err != nil {
					ch <- prometheus.NewInvalidMetric(c.echonetMetrics, err)
					return
				}
			}

			{
				value, time, err := echonet.CalcCumulativePower()
				if err != nil {
					ch <- prometheus.NewInvalidMetric(c.echonetMetrics, err)
					return
				}

				ch <- prometheus.NewMetricWithTimestamp(
					time,
					prometheus.MustNewConstMetric(
						c.echonetMetrics,
						prometheus.GaugeValue,
						value,
						*appliance.Device.Name,
						"cumulative_power",
					),
				)
			}
			{
				value, time, err := echonet.CalcInstantaneousPower()
				if err != nil {
					ch <- prometheus.NewInvalidMetric(c.echonetMetrics, err)
					return
				}

				ch <- prometheus.NewMetricWithTimestamp(
					time,
					prometheus.MustNewConstMetric(
						c.echonetMetrics,
						prometheus.GaugeValue,
						value,
						*appliance.Device.Name,
						"instantaneous_power",
					),
				)
			}
		}
	}
}

func (c *Collector) addAuthHeader(ctx context.Context, req *http.Request) error {
	req.Header.Set("Authorization", "Bearer "+c.natureToken)
	return nil
}
