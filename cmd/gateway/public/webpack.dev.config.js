'use strict';

var path = require('path');
var webpack = require('webpack');
var HtmlWebpackPlugin = require('html-webpack-plugin')
var ExtractTextPlugin = require('extract-text-webpack-plugin')
var CopyWebpackPlugin = require('copy-webpack-plugin');
/**
 * Module dependencies
 */
module.exports = {
    cache: false,

    entry: {
        angular: [
            './node_modules/angular/angular',
            './node_modules/angular-route/angular-route',
            './node_modules/angular-animate/angular-animate',
            './node_modules/angular-aria/angular-aria'
        ],
        bootstrap: [
            './node_modules/bootstrap/dist/js/bootstrap',
            './node_modules/angular-ui-bootstrap/dist/ui-bootstrap',
            './node_modules/angular-ui-bootstrap/dist/ui-bootstrap-tpls'
        ],
        main: __dirname + '/src/entry',
    },

    output: {
        path: 'dist/',
        publicPath: 'dist/',
        filename: '[name].bundle-[chunkhash].js'
    },

    plugins: [
        new webpack.optimize.CommonsChunkPlugin({
          names: [ 'bootstrap', 'angular']
        }),
        new HtmlWebpackPlugin({
            template: 'ejs/index.html',
            inject: 'body',
            filename: '../index.html',
            chunks: ["angular", 'bootstrap', "main"]
        }),
        new CopyWebpackPlugin([
          { from: 'node_modules/bootstrap/dist/css/bootstrap.min.css', to: 'bootstrap/css' },
          { from: 'node_modules/bootstrap/dist/css/bootstrap.min.css.map', to: 'bootstrap/css' },
          { from: 'node_modules/bootstrap/dist/css/bootstrap.min.css.map', to: 'bootstrap/css' },
          { from: 'node_modules/jquery/dist/jquery.min.js', to: './' }
        ])
      ]
    //watch:true
};