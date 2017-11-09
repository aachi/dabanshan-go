'use strict';
require(['angular', './controllers', 'angular-route'], 
  function(angular, controllers) {
    angular.module('app', ['ngRoute', 'ui.bootstrap'])
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
          .otherwise({redirectTo: '/login'});
    }])
    angular.bootstrap(document, ['app']); 
});