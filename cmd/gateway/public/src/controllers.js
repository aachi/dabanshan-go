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

    controllers.OrdersMgrCtrl = function($scope, $rootScope) {
        
    }
    controllers.OrdersMgrCtrl.$inject = ['$scope', '$rootScope'];
    
    controllers.ProductsMgrCtrl = function($scope, $rootScope, $q, $location, $uibModal) {
        $scope.showModal = function () {
            var modalInstance = $uibModal.open({
                templateUrl: '../components/productModal.html',
                controller: ['$scope', '$uibModal', controllers.NewProductCtrl],
                size: 'lg',
                resolve: {
                    
                }
            });
            return modalInstance;
        }

    }
    controllers.ProductsMgrCtrl.$inject = ['$scope','$rootScope','$q','$location', '$uibModal'];


    controllers.NewProductCtrl = function($scope, $rootScope) {

    }
    controllers.NewProductCtrl.$inject = ['$scope', '$rootScope'];
    
    return controllers;
});