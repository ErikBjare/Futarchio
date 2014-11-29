var gulp = require('gulp');

/* Tooling */
var concat = require('gulp-concat');
var del = require('del');
var rename = require('gulp-rename');
var flatten = require('gulp-flatten');
var gulpFilter = require('gulp-filter');

/* Minifying stuff */
var ngAnnotate = require('gulp-ng-annotate');
var uglify = require('gulp-uglify');
var htmlmin = require('gulp-htmlmin');
var imagemin = require('gulp-imagemin');
var sourcemaps = require('gulp-sourcemaps');

/* Stylesheet stuff */
var sass = require('gulp-sass');

/* Bower compatibility */
var mainBowerFiles = require('main-bower-files');


var root_dir = 'src';
var src_path = root_dir+'/site/src';
var dist_path = root_dir+'/site/dist';
var paths = {
  scripts: src_path+'/scripts/**/*.js',
  go: root_dir+'**/*.go',
  images: src_path+'/img/**/*',
  html: src_path+'/**/*.html',
  robots: src_path+'/robots.txt',
  stylesheets_main: src_path+'/stylesheets/style.scss',
  stylesheets_all: src_path+'/stylesheets/*.scss'
};

// Not all tasks need to use streams
// A gulpfile is just another node program and you can use all packages available on npm
gulp.task('clean', function(cb) {
  // You can use multiple globbing patterns as you would with `gulp.src`
  del([dist_path], cb);
});


// Minify and copy all JavaScript (except vendor scripts)
// with sourcemaps all the way down
gulp.task('scripts', [], function() {
  return gulp.src(paths.scripts)
    .pipe(sourcemaps.init())
      .pipe(ngAnnotate())
      .pipe(uglify())
      .pipe(concat('all.min.js'))
    .pipe(sourcemaps.write())
    .pipe(gulp.dest(dist_path));
});

// Copy and optimize all static images
gulp.task('images', [], function() {
  return gulp.src(paths.images)
    // Pass in options to the task
    .pipe(imagemin({optimizationLevel: 5}))
    .pipe(gulp.dest(dist_path+'/img'));
});

// Generate CSS from SASS
gulp.task('stylesheets', [], function() {
  return gulp.src(paths.stylesheets_main)
    .pipe(sass())
    .pipe(gulp.dest(dist_path));
});

gulp.task('html', [], function() {
  return gulp.src(paths.html)
    .pipe(htmlmin({collapseWhitespace:true}))
    .pipe(gulp.dest(dist_path));
});

gulp.task('robots', [], function() {
  return gulp.src(paths.robots)
    .pipe(gulp.dest(dist_path));
});

// Rerun the task when a file changes
gulp.task('watch', function() {
  gulp.watch(paths.scripts, ['scripts']);
  gulp.watch(paths.images, ['images']);
  gulp.watch(paths.stylesheets_all, ['stylesheets']);
  gulp.watch(paths.html, ['html']);
  gulp.watch(["bower.json"], ['libs']);
});

// grab libraries files from bower_components, minify and send to dist
gulp.task('libs', function() {

    var jsFilter = gulpFilter('*.js');
    var cssFilter = gulpFilter('*.css');
    var fontFilter = gulpFilter(['*.eot', '*.woff', '*.svg', '*.ttf']);

    return gulp.src(mainBowerFiles())

    // grab vendor js files from bower_components, minify and send to dist
    .pipe(jsFilter)
    .pipe(gulp.dest(dist_path + '/scripts'))
    .pipe(uglify())
    .pipe(rename({
        suffix: ".min"
    }))
    .pipe(gulp.dest(dist_path + '/scripts'))
    .pipe(jsFilter.restore())

    // grab vendor css files from bower_components, minify and send to dist
    .pipe(cssFilter)
    .pipe(gulp.dest(dist_path + '/stylesheets'))
    // TODO: .pipe(minifycss())
    .pipe(rename({
        suffix: ".min"
    }))
    .pipe(gulp.dest(dist_path + '/stylesheets'))
    .pipe(cssFilter.restore())

    // grab vendor font files from bower_components and send to dist
    .pipe(fontFilter)
    .pipe(flatten())
    .pipe(gulp.dest(dist_path + '/fonts'));
});

// The default task (called when you run `gulp` from cli)
gulp.task('build', ['libs', 'scripts', 'images', 'stylesheets', 'html', 'robots']);
gulp.task('default', ['build', 'watch']);
