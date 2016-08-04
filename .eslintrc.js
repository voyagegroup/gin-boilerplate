module.exports = {
    extends: [
        "eslint:recommended",
        "plugin:react/recommended"
    ],
    parser: 'babel-eslint',
    env: {
        es6: true,
        browser: true
    },
    plugins: [
        "react"
    ],
    ecmaFeatures: {
        jsx: true,
        modules: true
    },
    globals: {},
    rules: {
        'no-console': 0,
        'react/prop-types': 1
    }
};
