local grafana = import 'grafonnet-lib/grafonnet/grafana.libsonnet';
local dashboard = grafana.dashboard;
local row = grafana.row;
local singlestat = grafana.singlestat;
local prometheus = grafana.prometheus;
local template = grafana.template;

dashboard.new(
  'Keptn',
  schemaVersion=16,
  tags=['keptn'],
)
.addTemplate(
  grafana.template.datasource(
    'PROMETHEUS_DS',
    'prometheus',
    'Prometheus',
    hide='label',
  )
)
.addTemplate(
  template.new(
    'stage',
    '$PROMETHEUS_DS',
    'label_values(dev, staging)',
    label='Stage',
    refresh='time'
  )
)
.addPanel(
  grafana.graphPanel.new(
    'Response Time',
    datasource='$datasource',
    span=3,
  )
  .addTarget(
    grafana.prometheus.target(
      'histogram_quantile(0.9, sum by(le) (rate(http_response_time_milliseconds_bucket{job="carts-sockshop-production"}[$range_duration])))',
      legendFormat='response time p90'
    )
  ), gridPos={
    x:0,
    y:0,
    w:100,
    h: 3,
  },
)
