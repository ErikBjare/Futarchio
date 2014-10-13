var gulp = require('gulp');
var concat = require('gulp-concat');
var uglify = require('gulp-uglify');
var imagemin = require('gulp-imagemin');
var sourcemaps = require('gulp-sourcemaps');
var sass = require('gulp-sass');
var del = require('del');
var ngAnnotate = require('gulp-ng-annotate');
var htmlmin = require('gulp-htmlmin');

var src_path = 'src/site/src';
var dist_path = 'src/site/dist';
var paths = {
  scripts: src_path+'/scripts/**/*.js',
  images: src_path+'/img/**/*',
  html: src_path+'/**/*.html',
  robots: src_path+'/robots.txt',
  stylesheets: src_path+'/stylesheets/style.scss'
};

// Not all tasks need to use streams
// A gulpfile is just another node program and you can use all packages available on npm
gulp.task('clean', function(cb) {
  // You can use multiple globbing patterns as you would with `gulp.src`
  del([dist_path], cb);
});

gulp.task('scripts', [], function() {
  // Minify and copy all JavaScript (except vendor scripts)
  // with sourcemaps all the way down
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
  return gulp.src(paths.stylesheets)
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
  gulp.watch(paths.stylesheets, ['stylesheets']);
  gulp.watch(paths.html, ['html']);
});

// The default task (called when you run `gulp` from cli)
gulp.task('default', ['watch', 'scripts', 'images', 'stylesheets', 'html', 'robots']);
