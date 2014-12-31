(function () {

    'use strict';

    require('angular/angular');
    require('angular-bootstrap/ui-bootstrap-tpls');
    require('angular-cookies/angular-cookies');
    require('angular-ui-router/release/angular-ui-router');
    //require('bootstrap');
    //require('font-awesome');
    //require('jquery');
    //require('rdash-ui');
    var alertCtrl = require('./controllers/alert-ctrl');
    var masterCtrl = require('./controllers/master-ctrl.js');
    var rdLoading = require('./directives/loading');
    var rdWidget = require('./directives/widget');
    var rdWidgetBody = require('./directives/widget-body');
    var rdWidgetFooter = require('./directives/widget-footer');
    var rdWidgetTitle = require('./directives/widget-header');

    angular.module('RDash', ['ui.bootstrap', 'ui.router', 'ngCookies'])
        .config(['$stateProvider', '$urlRouterProvider',
            function($stateProvider, $urlRouterProvider) {

                // For unmatched routes
                $urlRouterProvider.otherwise('/');

                // Application routes
                $stateProvider
                    .state('index', {
                        url: '/',
                        templateUrl: 'templates/dashboard.html'
                    })
                    .state('tables', {
                        url: '/tables',
                        templateUrl: 'templates/tables.html'
                    });
            }
        ])
        .controller('AlertsCtrl', ['$scope', alertCtrl])
        .controller('MasterCtrl', ['$scope', '$cookieStore', masterCtrl])
        .directive('rdLoading', rdLoading)
        .directive('rdWidget', rdWidget)
        .directive('rdWidgetBody', rdWidgetBody)
        .directive('rdWidgetFooter', rdWidgetFooter)
        .directive('rdWidgetHeader', rdWidgetTitle);
}());
