import gulp from 'gulp'
import gutil from 'gulp-util'
import runSequence from 'run-sequence'
import webpackStream from 'webpack-stream'
import rename from 'gulp-rename'
import webpackServer from 'webpack-dev-server'
import webpack from 'webpack'
import { readFileSync } from 'fs'

const hashOptions = {
    algorithm: 'md5',
    hashLength: 10,
    template: '<%= name %>.<%= hash %><%= ext %>'
};

const host = argv.host || 'localhost'
const port = argv.port || 3030

const productionLoaders = {
    loaders: [
        {
            test: /\.jsx?$/,
            //loader: 'babel-loader',
            exclude: /node_modules/,
            //query: babelConfig
        },
    ]
}

const devLoaders = {
    loaders: [
        {
            test: /\.jsx?$/,
            loaders: [
                //'babel-loader?' + JSON.stringify(babelConfig)
            ],
            exclude: /node_modules/,
        },
    ]
}

const webpackConfig = {
    cache: true,
    context: __dirname + '/src/js',
    entry: {
        'main': './containers/App/App.jsx',
        'admin': './components/Admin/Admin.jsx',
    },
    output: {
        path: __dirname + '/dist/js',
        filename: '[name].js',
        publicPath: 'http://'+host+':'+port+'/dist'
    },
    plugins: [],
    resolve: {
        root: [
            __dirname + '/app/js',
        ],
        extensions: ['', '.js', '.jsx'],
    },
}

gulp.task('clean', cb => {
    del(['dist'], cb)
})

gulp.task('js:webpack', cb => {

    webpackConfig.module = productionLoaders
    webpackConfig.plugins.push(
        new webpack.DefinePlugin({
            __CLIENT__: true,
            __SERVER__: false,
            __DEVELOPMENT__: false,
            __DEVTOOLS__: false,  // <-------- DISABLE redux-devtools HERE
            'process.env': {
                'NODE_ENV': JSON.stringify('production')
            }
        })
    )

    return gulp.src(paths.browser_js)
        .pipe(webpackStream(webpackConfig))
        .pipe(gulp.dest('dist/js'))
})

gulp.task('index:build', cb => {
    return gulp.src('static/index.html.tmpl')
        .pipe(mustache({
            dev: process.env.NODE_ENV !== 'production',
            devserverHost: host,
            devserverPort: port,
        }))
        .pipe(rename('index.html'))
        .pipe(gulp.dest('static'))
})

const devCompiler = webpack(webpackConfig);

gulp.task("webpack-dev-server", function(cb) {
    // modify some webpack config options
    webpackConfig.module = devLoaders;
    webpackConfig.devtool = 'eval',
    webpackConfig.debug = true;

    for (var i in webpackConfig.entry) {
        var originalEntry = webpackConfig.entry[i];
        webpackConfig.entry[i] = [
            'webpack-dev-server/client?http://'+host+':'+port,
            'webpack/hot/only-dev-server',
            originalEntry
        ];
    }

    webpackConfig.plugins.push(new webpack.HotModuleReplacementPlugin());
    webpackConfig.plugins.push(new webpack.NoErrorsPlugin());

    webpackConfig.plugins.push(
        new webpack.DefinePlugin({
            __CLIENT__: true,
            __SERVER__: false,
            __DEVELOPMENT__: true,
            __DEVTOOLS__: true,
        })
    );

    // Start a webpack-dev-server
    new WebpackServer(webpack(webpackConfig), {
        contentBase: "dist",
        publicPath: 'http://'+host+':'+port+'/dist',
        hot: true,
        headers: { 'Access-Control-Allow-Origin': '*' },
        stats: {
          colors: true,
          progress: true
        }
    }).listen(port, host, function(err) {
        if(err) throw new gutil.PluginError("webpack-dev-server", err);
        gutil.log("[webpack-dev-server]", 'http://'+host+':'+port+'/webpack-dev-server/index.html');
    });

    var gracefulShutdown = function() {
        process.exit()
    }

    // listen for INT signal e.g. Ctrl-C
    process.on('SIGINT', gracefulShutdown);
});

gulp.task('dist', ['js:dist', 'index:build'])
gulp.task('default', ['webpack-dev-server'])
