import { DataSourcePlugin } from '@grafana/data';
import { AutotaskDatasource } from './datasource';
import { ConfigEditor } from './components/ConfigEditor';
import { QueryEditor } from './components/QueryEditor';
import { AutotaskQuery, AutotaskDatasourceOptions, AutotaskSecureJsonData } from './types';

export const plugin = new DataSourcePlugin<
  AutotaskDatasource,
  AutotaskQuery,
  AutotaskDatasourceOptions,
  AutotaskSecureJsonData
>(AutotaskDatasource)
  .setConfigEditor(ConfigEditor)
  .setQueryEditor(QueryEditor);
