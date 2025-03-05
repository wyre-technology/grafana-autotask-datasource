import { DataSourcePlugin } from '@grafana/data';
import { ConfigEditor } from './components/ConfigEditor';
import { AutotaskDataSource } from './datasource';
import { AutotaskDataSourceOptions, AutotaskQuery } from './types';

export const plugin = new DataSourcePlugin<AutotaskDataSource, AutotaskQuery, AutotaskDataSourceOptions>(
  AutotaskDataSource
)
  .setConfigEditor(ConfigEditor)
  .setMetadata({
    id: 'wyretech-autotask-datasource',
    name: 'Autotask',
    description: 'Query data from Autotask PSA',
    info: {
      author: {
        name: 'Wyretech',
      },
      links: [
        {
          name: 'GitHub',
          url: 'https://github.com/wyretech/autotask-datasource',
        },
      ],
      version: '1.0.0',
      updated: '2024-03-05',
    },
  });
