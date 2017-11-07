'use strict';

var path = require('path');
var webpack = require('webpack');
var HtmlWebpackPlugin = require('html-webpack-plugin')
var ExtractTextPlugin = require('extract-text-webpack-plugin')
    
/**
 * Module dependencies
 */
module.exports = {
    cache: false,
    entry: {
        'app': __dirname + '/src/entry',
    },

    output: {
        path: 'dist/',
        publicPath: 'dist/',
        filename: '[name].js',
        chunkFilename: '[chunkhash].js'
    },
    //watch:true
};