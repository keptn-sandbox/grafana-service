'use strict';

var grafana = require('grafana-dash-gen');
var Row = grafana.Row;
var Dashboard = grafana.Dashboard;
var Panels = grafana.Panels;
var Target = grafana.Target;
var Templates = grafana.Templates;

// For grafana v1, the URL should look something like:
//
//   https://your.grafana.com/elasticsearch/grafana-dash/dashboard/
//
// Bascially, grafana v1 used elastic search as its backend, but grafana v2
// has its own backend. Because of this, the URL for grafana v2 should look
// something like this:
//
//   https://your.grafanahost.com/grafana2/api/dashboards/db/

// grafana.configure({
//     url: 'http://35.222.54.134/grafana2/api/dashboards/db/',
//     cookie: 'auth-openid=eyJrIjoiZVNGMms4eE1NYTVMRTJTc3dneVYxVDdDYUt2RFcyQ2YiLCJuIjoiZGFzaGJvYXJkLWdlbiIsImlkIjoxfQ==',
//     headers: {'x': 'y'}
// });

grafana.configure({
    url: 'http://35.222.54.134/api/dashboards/db',
    authorization: 'Bearer eyJrIjoiZVNGMms4eE1NYTVMRTJTc3dneVYxVDdDYUt2RFcyQ2YiLCJuIjoiZGFzaGJvYXJkLWdlbiIsImlkIjoxfQ==',
});

// Dashboard Constants
var TITLE = 'TEST API dashboard';
var TAGS = ['myapp', 'platform'];
var TEMPLATING = [{
    name: 'dc',
    options: ['dc1', 'dc2']
}, {
    name: 'smoothing',
    options: ['30min', '10min', '5min', '2min', '1min']
}];
var ANNOTATIONS = [{
    name: 'Deploy',
    target: 'stats.$dc.production.deploy'
}];
var REFRESH = '1m';

// Target prefixes
var SERVER_PREFIX = 'servers.app*-$dc.myapp.';
var COUNT_PREFIX = 'stats.$dc.counts.myapp.';

function generateDashboard() {
    // Rows
    var volumeRow = new Row({
        title: 'Request Volume'
    });
    var systemRow = new Row({
        title: 'System / OS'
    });

    // Panels: request volume
    var rpsGraphPanel = new Panels.Graph({
        title: 'req/sec',
        span: 8,
        targets: [
            new Target(COUNT_PREFIX + 'statusCode.*').
                        transformNull(0).
                        sum().
                        hitcount('1seconds').
                        scale(0.1).
                        alias('rps')
        ]
    });
    var rpsStatPanel = new Panels.SingleStat({
        title: 'Current Request Volume',
        postfix: 'req/sec',
        span: 4,
        targets: [
            new Target(COUNT_PREFIX + 'statusCode.*').
                    sum().
                    scale(0.1)
        ]
    });

    var favDashboardList = new Panels.DashboardList({
        title: 'My Favorite Dashboard',
        span: 4,
        mode: 'search',
        query: 'dashboard list'
    });

    // Panels: system health
    var cpuGraph = new Panels.Graph({
        title: 'CPU',
        span: 4,
        targets: [
            new Target(SERVER_PREFIX + 'cpu.user').
                nonNegativeDerivative().
                scale(1 / 60).
                scale(100).
                averageSeries().
                alias('avg'),
            new Target(SERVER_PREFIX + 'cpu.user').
                nonNegativeDerivative().
                scale(1 / 60).
                scale(100).
                percentileOfSeries(95, false).
                alias('p95')
        ]
    });
    var rssGraph = new Panels.Graph({
        title: 'Memory',
        span: 4,
        targets: [
            new Target(SERVER_PREFIX + 'memory.rss').
                averageSeries().
                alias('rss')
        ]
    });
    var fdsGraph = new Panels.Graph({
        title: 'FDs',
        span: 4,
        targets: [
            new Target(SERVER_PREFIX + 'fds').
                averageSeries().
                movingAverage('10min').
                alias('moving avg')
        ]
    });

    // Dashboard
    var dashboard = new Dashboard({
      title: 'Api dashboard'
    })
    var requestVolume = new Panels.SingleStat({
      title: 'Current Request Volume',
      postfix: 'req/sec',
      targets: [
        new Target('stats.$dc.counts').
            sum().scale(0.1)
      ],
      row: new Row(),
      dashboard: dashboard
    });
    // Finish
    console.log(dashboard.generate());
    //grafana.publish(dashboard);
}

module.exports = {
    generate: generateDashboard
};

if (require.main === module) {
    generateDashboard();
}
