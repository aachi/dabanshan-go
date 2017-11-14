'use strict';

// require(['angular', './controllers', './services', 'bootstrap'],
//     function (angular, controllers) {
//         angular.module('dbsApp', ['app.services', 'ngRoute', 'ui.router', 'ui.bootstrap'])
//             .config(['$stateProvider', '$urlRouterProvider', function ($stateProvider, $urlRouterProvider, $httpProvider) {
//                 $urlRouterProvider.otherwise('/login');
//                 $stateProvider
//                     .state('site.login', { url: '/login', controllerAs: 'login', templateUrl: 'views/login.html?' + +new Date(), controller: controllers.HomeCtrl })
//                     .state('site.explore', { url: '/explore', controllerAs: 'explore', templateUrl: 'views/explore.html?' + +new Date(), controller: controllers.ExploreCtrl })
//                     .state('mgr.tenants', { url: '/tenants', controllerAs: 'tenants', templateUrl: 'views/tenants.html?' + +new Date(), controller: controllers.TenantsCtrl })
//                     .state('mgr.dashboard', { url: '/dashboard', controllerAs: 'dashboard', templateUrl: 'views/dashboard.html?' + +new Date(), controller: controllers.DashboardCtrl })
//                     .state('mgr.orders', { url: '/mgr/orders', templateUrl: 'views/m-orders.html?' + +new Date(), controller: controllers.OrdersMgrCtrl })
//                     .state('mgr.products', { url: '/mgr/products', templateUrl: 'views/m-products.html?' + +new Date(), controller: controllers.ProductsMgrCtrl })
//             }])
//         angular.bootstrap(document, ['dbsApp']);
//     });

var controllers = require('./controllers')
var services = require('./services')
var bootstrap = require('bootstrap')
var routerApp = angular.module('dbsApp', ['ui.router', 'ui.bootstrap', 'app.services']);

routerApp.config(function ($stateProvider, $urlRouterProvider) {
    $urlRouterProvider.otherwise('/login');
    $stateProvider
        .state('login', { url: '/login', templateUrl: 'views/login.html', controller: controllers.HomeCtrl })
        .state('explore', { url: '/explore', controllerAs: 'explore', templateUrl: 'views/explore.html?' + +new Date(), controller: controllers.ExploreCtrl })
        .state('tenants', { url: '/tenants', controllerAs: 'tenants', templateUrl: 'views/tenants.html?' + +new Date(), controller: controllers.TenantsCtrl })
        .state('dashboard', { url: '/dashboard', controllerAs: 'dashboard', templateUrl: 'views/dashboard.html?' + +new Date(), controller: controllers.DashboardCtrl })
        .state('dashboard.orders', { url: '/orders', templateUrl: 'views/m-orders.html?' + +new Date(), controller: controllers.OrdersMgrCtrl })
        .state('dashboard.products', { url: '/products', templateUrl: 'views/m-products.html?' + +new Date(), controller: controllers.ProductsMgrCtrl })
});