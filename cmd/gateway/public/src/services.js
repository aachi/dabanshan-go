'use strict';
define(['angular'], function (angular) {
    /* Services */
    angular.module('app.services', [])
        .factory('Config', [function () {
            var Config = {
                url: "http://localhost:8000/api/v1/"
            };
            return Config;
        }])
        // with users
        .factory('UserService', ['$http', 'Config', function ($http, Config) {
            var User = {
                type: 'users',
                userOBJ: {},
                username: "",
                password: "",
                token: "",
                setUsername: function (inputUsername) {
                    User.username = inputUsername;
                },

                getUserName: function () {
                    return this.username;
                },

                setToken: function (token) {
                    User.token = token;
                },

                getToken: function () {
                    return this.token;
                },

                getLoggedUserName: function () {
                    return localStorage.getItem("username");
                },

                setPassword: function (inputPassword) {
                    User.password = inputPassword;
                },
                getPassword: function () {
                    if (this.password != "") {
                        return this.password;
                    }
                    else {
                        return localStorage.getItem("password");
                    }
                },
                storeUserLocally: function (obj) {
                    if (typeof (Storage) !== "undefined") {
                        localStorage.setItem("username", this.getUserName());
                        localStorage.setItem("password", this.getPassword());
                        localStorage.setItem("token", this.getToken());
                    } else {
                        console.log('no local storage available');
                    }
                },

                getUserOBJ: function () {
                    return $this.userOBJ;
                },
                login: function (inputUsername, inputPassword, callback) {
                    var headers = { 'Content-Type': 'application/json' };
                    var jsonObject = angular.toJson({ "username": inputUsername, "password": inputPassword });
                    $http.post(Config.url + this.type + '/login', jsonObject, { headers: headers })
                        .then(function (response) {
                            User.setUsername(response.user.username);
                            User.setToken(response.token);
                            User.storeUserLocally(response.user);
                            callback(response);
                        })
                        .catch(function (err) {
                            callback(err);
                        })
                },
                logout: function () {

                }
            }
            return User;
        }])
        // with products
        .factory('ProductService', ['$http', 'Config', function ($http, Config) {
            var Product = {
                type: 'products',
                getProductsByTenant: function (params, callback) {

                    return ""
                }
            }
            return Product
        }])
        // with orders
        .factory('OrderService', ['$http', 'Config', function ($http, Config) {
            var Order = {
                type: 'orders',
                getOrdersByTenant: function (tenantID, callback) {
                    var headers = { 'Content-Type': 'application/json' };
                    var jsonObject = angular.toJson({ "tenantID": tenantID});
                    $http.get(Config.url + this.type + '/', jsonObject, { headers: headers })
                        .success(function (response) {
                            callback(response);
                        })
                        .error(function (err) {
                            callback(err);
                        })
                }
            }
            return Order
        }])
});