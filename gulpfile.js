var config       = require('./gulp.config.json');
var gulp         = require('gulp');
var browserSync  = require('browser-sync');
var sass         = require('gulp-sass');
var concat       = require('gulp-concat');
var postcss      = require('gulp-postcss');
var autoprefixer = require('autoprefixer');
var postcssFlexbugsFixes = require('postcss-flexbugs-fixes');

gulp.task('sass', ['copy-bower-files'], function() {
    return gulp.src(config.styles.files)
        .pipe(postcss([ autoprefixer({ browsers: ['last 2 versions'] }), postcssFlexbugsFixes ]))
        .pipe(sass())
        .pipe(gulp.dest(config.styles.dest));
});

gulp.task('scripts', function () {
    return gulp.src(config.scripts.files)
        .pipe(concat(config.scripts.filename))
        .pipe(gulp.dest(config.scripts.dest));
});

gulp.task('copy-bower-files', function () {
    return gulp.src(config.bower.files)
        .pipe(gulp.dest(config.bower.dest));
});

gulp.task('browser-sync', function() {
    browserSync.init(config.browsersync);
});

gulp.task('watch', function() {
    gulp.watch(config.styles.files , ['sass']);
    gulp.watch(config.scripts.files, ['scripts']);
});

gulp.task('default', ['sass' , 'scripts']);
