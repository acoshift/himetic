const gulp = require('gulp')
const sass = require('gulp-sass')
const rev = require('gulp-rev')
const uglify = require('gulp-uglify')
const rename = require('gulp-rename')
const rollup = require('gulp-rollup')
const resolve = require('rollup-plugin-node-resolve')
const babel = require('rollup-plugin-babel')
const map = require('map-stream')
const fs = require('fs')
const path = require('path')

const style = './style/main.@(sass|scss)'
const script = './src/main.js'
const assets = './assets'
const static = './static.yaml'

function writeFilename(name) {
  return function (file, cb) {
    const p = path.basename(file.path)
    fs.appendFileSync(static, `${name}: ${p}\n`)
    cb(null, file)
  }
}

gulp.task('default', ['clear', 'style', 'script'])

gulp.task('clear', function (cb) {
  fs.unlink(static, () => cb())
})

gulp.task('style', () => gulp
  .src(style)
  .pipe(sass.sync({
    outputStyle: 'compressed',
    includePaths: './node_modules',
  }).on('error', sass.logError))
  .pipe(rename('app.css'))
  .pipe(rev())
  .pipe(map(writeFilename('app.css')))
  .pipe(gulp.dest(assets))
)

gulp.task('script', () => gulp
  .src(script)
  .pipe(rollup({
    input: script,
    format: 'iife',
    plugins: [
      resolve(),
      babel({
        babelrc: false,
        presets: [
          ['env', {
            'modules': false
          }]
        ]
      })
    ]
  }))
  .pipe(uglify())
  .pipe(rename('app.js'))
  .pipe(rev())
  .pipe(map(writeFilename('app.js')))
  .pipe(gulp.dest(assets))
)
