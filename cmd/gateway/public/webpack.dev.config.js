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
            './node_modules/angular-route/angular-route'
        ],
        bootstrap: './node_modules/angular-ui-bootstrap/dist/ui-bootstrap-tpls',
        main: __dirname + '/src/entry',
    },

    output: {
        path: 'dist/',
        publicPath: 'dist/',
        filename: '[name].bundle.js'
    },

    // module:{
    //     rules: [{
    //         test: /\.css$/,
    //         use: ExtractTextPlugin.extract({
    //             fallback: "style-loader",
    //             use: [{
    //                 loader: "css-loader"
    //             },],
    //         })
    //     }]
    // },

    plugins: [
        // new ExtractTextPlugin('styles/main.css'),
        // new webpack.optimize.UglifyJsPlugin({
        //   sourceMap: true
        // }),
        new webpack.optimize.CommonsChunkPlugin({
          names: [ 'bootstrap', 'angular' ]
        }),
        new CopyWebpackPlugin([
          { from: 'node_modules/bootstrap/dist/css/bootstrap.min.css', to: 'bootstrap/css' },
          { from: 'node_modules/bootstrap/dist/fonts', to: 'bootstrap/fonts' },
        ])
      ]
    //watch:true
};