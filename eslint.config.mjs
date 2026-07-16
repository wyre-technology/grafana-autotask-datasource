// Flat-config migration (2026-07-16): @grafana/eslint-config@10 + eslint@9 dropped
// support for the legacy .eslintrc format this project used previously
// (.config/.eslintrc). See https://eslint.org/docs/latest/use/configure/migration-guide.
import grafanaConfig from '@grafana/eslint-config';

export default [
  {
    ignores: [
      'dist/',
      'node_modules/',
      'coverage/',
      'work/',
      'artifacts/',
      'ci/',
      'playwright-report/',
      'test-results/',
      'blob-report/',
      '.config/',
    ],
  },
  {
    // ESLint 9's default file discovery only covers .js/.mjs/.cjs; widen it
    // so the TS/TSX source tree actually gets visited.
    files: ['**/*.{js,jsx,ts,tsx,mjs,cjs,mts,cts}'],
  },
  ...grafanaConfig,
  {
    rules: {
      'react/prop-types': 'off',
    },
  },
  {
    files: ['tests/**/*'],
    rules: {
      'react-hooks/rules-of-hooks': 'off',
    },
  },
];
