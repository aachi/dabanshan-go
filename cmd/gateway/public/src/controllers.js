/*global define */

'use strict';

define(function() {
    var controllers = {};

    controllers.HomeCtrl = function($scope, $rootScope, $location, UserService) {
        $scope.login = function() {
            UserService.login($scope.username, $scope.password, function(resp){
                $location.path("/explore")
            })
        }
    }
    controllers.HomeCtrl.$inject = ['$scope', '$rootScope', '$location', 'UserService'];
    
    controllers.ExploreCtrl = function($scope, $rootScope) {
        
    }
    controllers.ExploreCtrl.$inject = ['$scope', '$rootScope'];
 
    controllers.TenantsCtrl = function($scope, $rootScope) {
        
    }
    controllers.TenantsCtrl.$inject = ['$scope', '$rootScope'];

    controllers.DashboardCtrl = function($scope, $rootScope) {
        
    }
    controllers.DashboardCtrl.$inject = ['$scope', '$rootScope'];


    return controllers;
});