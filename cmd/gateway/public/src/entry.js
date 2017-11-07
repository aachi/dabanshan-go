'use strict';

require(['angular', './controllers', 'angular-route'], 
  function(angular, controllers) {
    angular.module('app', ['ngRoute'])
      .config(['$routeProvider', '$httpProvider', function ($routeProvider, $httpProvider) {
        $routeProvider.otherwise({redirectTo: '/home'});
        $routeProvider.when('/home', {templateUrl: 'views/home.html', controller: controllers.HomeCtrl})
    }])
    angular.bootstrap(document, ['app']); 
});