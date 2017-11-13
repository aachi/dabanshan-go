'use strict';
require(['angular', './controllers', './services', 'angular-route', 'angular-animate', 'bootstrap'], 
  function(angular, controllers) {
    angular.module('app', ['app.services', 'ngRoute', 'ui.bootstrap'])
        .run(function($rootScope, $templateCache) {  
          $rootScope.$on('$routeChangeStart', function(event, next, current) {  
              if (typeof(current) !== 'undefined'){  
                  $templateCache.remove(current.templateUrl);  
              }  
          });  
      })
      .config(['$routeProvider', '$httpProvider', function ($routeProvider, $httpProvider) {
        $routeProvider
          .when('/login', {templateUrl: 'views/login.html?' + +new Date(), controller: controllers.HomeCtrl})
          .when('/explore', {templateUrl: 'views/explore.html?' + +new Date(), controller: controllers.ExploreCtrl})
          .when('/tenants', {templateUrl: 'views/tenants.html?' + +new Date(), controller: controllers.TenantsCtrl})
          .when('/dashboard', {templateUrl: 'views/dashboard.html?' + +new Date(), controller: controllers.DashboardCtrl})
          .when('/m-orders', {templateUrl: 'views/m-orders.html?' + +new Date(), controller: controllers.OrdersMgrCtrl})
          .when('/m-products', {templateUrl: 'views/m-products.html?' + +new Date(), controller: controllers.ProductsMgrCtrl})
          .otherwise({redirectTo: '/login'});
    }])
    angular.bootstrap(document, ['app']); 
});