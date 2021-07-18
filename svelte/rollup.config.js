import svelte from 'rollup-plugin-svelte'
import resolve from '@rollup/plugin-node-resolve'
import sveltePreprocess from 'svelte-preprocess'
import typescript from '@rollup/plugin-typescript'
import commonjs from '@rollup/plugin-commonjs'
import { terser } from 'rollup-plugin-terser'
import css from 'rollup-plugin-css-only'

const production = process.env.MIX_ENV === 'prod'

export default {
    input: 'src/main.ts',
    output: {
        sourcemap: false,
        format: 'iife',
        name: 'hangman',
        file: '../priv/static/app.js'
    },
    plugins: [
        svelte({
            preprocess: sveltePreprocess({ sourceMap: false }),
            compilerOptions: {
                dev: !production
            }
        }),
        css({ output: 'app.css' }),
        resolve({
            browser: true,
            dedupe: ['svelte']
        }),
        commonjs(),
        typescript({ sourceMap: false }),
        production && terser()
    ],
    watch: {
        clearScreen: false
    }
}