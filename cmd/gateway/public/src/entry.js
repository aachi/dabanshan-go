'use strict';
require(['angular', './controllers', 'angular-route'], 
  function(angular, controllers) {
    angular.module('app', ['ngRoute', 'ui.bootstrap'])
      .config(['$routeProvider', '$httpProvider', function ($routeProvider, $httpProvider) {
        $routeProvider.otherwise({redirectTo: '/login'});
        $routeProvider.when('/login', {templateUrl: 'views/login.html', controller: controllers.HomeCtrl})
    }])
    angular.bootstrap(document, ['app']); 
});