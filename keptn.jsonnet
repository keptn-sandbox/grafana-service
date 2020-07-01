local grafana = import 'grafonnet-lib/grafonnet/grafana.libsonnet';
local dashboard = grafana.dashboard;
local row = grafana.row;
local template = grafana.template;
local singlestat = grafana.singlestat;
local graphPanel = grafana.graphPanel;
local prometheus = grafana.prometheus;

local myprometheusDS = 'Prometheus';

local environment(env,i) = 
  grafana.text.new(
    title="Stage",
    content="# "+env,
    
  );


local response_time(job) = 
  graphPanel.new(
    title='Response time p90',
    datasource='$datasource',
  ).addTarget(
    prometheus.target(
      'histogram_quantile(0.9, sum by(le) (rate(http_response_time_milliseconds_bucket{job="'+job+'"}[1m])))',
      legendFormat='response time p90'
    )
  );

local throughput(job) = 
  graphPanel.new(
    title='Throughput',
    datasource='$datasource',
  ).addTarget(
    prometheus.target(
      'sum(rate(http_requests_total{job="'+job+'"}[1m]))',
      legendFormat='throughput'
    )
  );

local error_rate(job) = 
  graphPanel.new(
    title='Error rate',
    datasource='$datasource',
    min=0,
    max=1
  ).addTarget(
    prometheus.target(
      'sum(rate(http_requests_total{job="'+job+'",status!~"2.."}[1m]))/sum(rate(http_requests_total{job="'+job+'"}[1m]))',
      legendFormat='error rate',
    )
  );


local row(stage, job) = 
  grafana.row.new(
    title=stage,
    showTitle=true,
  );
//local envs = std.mapWithIndex(function(e,i) environment(e,i) {gridPos: { h:5,w:6, x:i*6, y:0} }, std.extVar('stages'));

local make_row(stage) = [
    row(stage, std.extVar('service')+"-"+std.extVar('project')+"-"+stage) { gridPos: { h:5,w:24, x:0, y:0} },
    response_time(std.extVar('service')+"-"+std.extVar('project')+"-"+stage) { gridPos: { h:10,w:6, x:0, y:0 } },
    throughput(std.extVar('service')+"-"+std.extVar('project')+"-"+stage)    { gridPos: { h:10,w:6, x:6, y:0 } },
    error_rate(std.extVar('service')+"-"+std.extVar('project')+"-"+stage)    { gridPos: { h:10,w:6, x:12, y:0 } },

];

////////////////////
// CREATE DASHBOARD
////////////////////
dashboard.new(
  title='Keptn',
  schemaVersion=16,
  tags=['keptn'],
  time_from='now-1h',
  refresh='30s',
  editable=true,
)
.addTemplate(
  grafana.template.datasource(
    'datasource',
    'prometheus',
    myprometheusDS,
    hide='label',
  )
)
.addPanels(

  std.flattenArrays(std.map(make_row, std.reverse(std.extVar('stages'))))

  // [
  //   row('production', std.extVar('service')+"-"+std.extVar('project')+"-production") { gridPos: { h:5,w:24, x:0, y:0} },

  //   response_time(std.extVar('service')+"-"+std.extVar('project')+"-production") { gridPos: { h:10,w:6, x:0, y:0 } },
  //   throughput(std.extVar('service')+"-"+std.extVar('project')+"-production")    { gridPos: { h:10,w:6, x:6, y:0 } },
  //   error_rate(std.extVar('service')+"-"+std.extVar('project')+"-production")    { gridPos: { h:10,w:6, x:12, y:0 } },

  // ]+


)