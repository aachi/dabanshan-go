var gulp = require('gulp'),
ngmin = require('gulp-ngmin'),
concat = require('gulp-concat'),
autoprefix = require('gulp-autoprefixer'),
minifyCSS = require('gulp-minify-css'),
rename = require('gulp-rename'),
less = require('gulp-less'),
gutil = require('gulp-util'),
filesize = require('gulp-filesize'),
watch = require('gulp-watch'),
uglify = require('gulp-uglify'),
clean = require('gulp-clean'),
jshint = require("gulp-jshint"),
browserify = require('browserify'),
streamify = require('gulp-streamify'),
buffer = require('gulp-buffer'),
useref = require('gulp-useref'),
filter = require('gulp-filter'),
source = require('vinyl-source-stream'),
sourcemaps = require('gulp-sourcemaps');

gulp.task('browserify', function() {
gulp.start('clean:js');
return browserify('./js/app.js', {entry: true, debug: true})
    .bundle()
    .pipe(source('app.js'))
    .pipe(streamify(uglify()))
    .pipe(gulp.dest('./js/dist/'))
    .pipe(buffer())
    .pipe(sourcemaps.init({loadMaps: true}))
    .pipe(sourcemaps.write('./'))
    .pipe(gulp.dest('./js/dist/'));
});

gulp.task('watch:js', function() {
watch({glob: "./js/**/*.js"}, function() {
    gulp.start("browserify");
});
});

gulp.task('clean:js', function() {
return gulp.src(['./js/dist/*.js', './js/dist/*.js.map'])
    .pipe(clean());
});