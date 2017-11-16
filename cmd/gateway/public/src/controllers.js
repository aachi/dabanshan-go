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

    controllers.OrdersMgrCtrl = function($scope, $rootScope, OrderService) {
        $scope.orders = []
        init = function() {
            OrderService.getOrdersByTenant("", function(res){
                console.log(res)
                $scope.orders = res.orders.Data
            }, function(err){
                console.log("err:", err)
            })
        }
        init();
    }
    controllers.OrdersMgrCtrl.$inject = ['$scope', '$rootScope', 'OrderService'];
    
    controllers.ProductsMgrCtrl = function($scope, $rootScope, $q, $location, $uibModal) {
        $scope.showModal = function () {
            var modalInstance = $uibModal.open({
                templateUrl: '../components/productModal.html?3',
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