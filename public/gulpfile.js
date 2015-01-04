var gulp = require('gulp'),
    usemin = require('gulp-usemin'),
    wrap = require('gulp-wrap'),
    connect = require('gulp-connect'),
    watch = require('gulp-watch'),
    minifyCss = require('gulp-minify-css'),
    minifyJs = require('gulp-uglify'),
    concat = require('gulp-concat'),
    less = require('gulp-less'),
    rename = require('gulp-rename'),
    browserify = require('gulp-browserify'),
    minifyHTML = require('gulp-minify-html');

var paths = {
    bundled: 'app/js/bundled.js',
    main: 'app/js/main.js',
    styles: 'app/less/**/*.*',
    images: 'app/img/**/*.*',
    templates: 'app/templates/**/*.html',
    index: 'app/index.html',
    bower_fonts: 'app/bower_components/**/*.{ttf,woff,eof,svg}',
};

//watch rdash
gulp.task('custom-css-rdash', function() {
    return gulp.src(['./app/bower_components/rdash-ui/dist/css/*.css', '!./app/bower_components/rdash-ui/dist/css/*.min.css'])
        .pipe(minifyCss({keepBreaks: true}))
        .pipe(rename({ suffix: '.min' }))
        .pipe(gulp.dest('./app/bower_components/rdash-ui/dist/css/'));
});


/**
* Handle bower bower_components from index
*/
gulp.task('usemin', function() {
    return gulp.src(paths.index)
        .pipe(usemin({
            js: [minifyJs(), 'concat'],
            css: [minifyCss({keepSpecialComments: 0}), 'concat'],
        }))
        .pipe(gulp.dest('dist/'));
});

/**
 * Copy assets
 */
gulp.task('build-assets', ['copy-bower_fonts']);

gulp.task('copy-bower_fonts', function() {
    return gulp.src(paths.bower_fonts)
        .pipe(rename({
            dirname: '/fonts'
        }))
        .pipe(gulp.dest('dist/lib'));
});

/**
 * Handle custom files
 */
gulp.task('build-custom', ['usemin', 'custom-images','browserify', 'custom-js', 'custom-less', 'custom-templates', 'custom-css-rdash']);

gulp.task('custom-images', function() {
    return gulp.src(paths.images)
        .pipe(gulp.dest('dist/img'));
});

gulp.task('custom-js', function() {
    return gulp.src(paths.bundled)
        .pipe(minifyJs())
        .pipe(rename({ suffix: '.min' }))
        .pipe(gulp.dest('dist/js'));
});

gulp.task('custom-less', function() {
    return gulp.src(paths.styles)
        .pipe(less())
        .pipe(gulp.dest('dist/css'));
});

gulp.task('custom-templates', function() {
    return gulp.src(paths.templates)
        .pipe(minifyHTML())
        .pipe(gulp.dest('dist/templates'));
});

/**
 * Watch custom files
 */
gulp.task('watch', function() {
    gulp.watch([paths.images], ['custom-images']);
    gulp.watch([paths.styles], ['custom-less']);
    gulp.watch([paths.main], ['browserify']);
    gulp.watch([paths.bundled], ['custom-js']);
    gulp.watch([paths.templates], ['custom-templates']);
    gulp.watch([paths.index], ['usemin']);
});

/**
 * Live reload server
 */
gulp.task('webserver', function() {
    connect.server({
        root: 'dist',
        livereload: true,
        port: 8888
    });
});

gulp.task('livereload', function() {
    gulp.src(['dist/**/*.*'])
        .pipe(watch())
        .pipe(connect.reload());
});

//browerify
gulp.task('browserify', function() {
    gulp.src(paths.main)
        .pipe(browserify({
            insertGlobals: true,
            debug: true
        }))
        .pipe(concat('bundled.js'))
        .pipe(gulp.dest('./app/js'))
});
/**
 * Gulp tasks
 */


gulp.task('build', ['build-assets', 'build-custom']);
gulp.task('default', ['build', 'webserver', 'livereload', 'watch']);