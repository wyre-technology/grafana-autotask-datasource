import {
  DataSourceApi,
  DataSourceInstanceSettings,
  DataSourcePluginMeta,
  PluginMetaInfo,
} from '@grafana/data';
import { AutotaskQuery, AutotaskDataSourceOptions } from './types';

export class AutotaskDataSource extends DataSourceApi<AutotaskQuery, AutotaskDataSourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<AutotaskDataSourceOptions>) {
    super(instanceSettings);
  }

  async query(options: any) {
    const { targets } = options;
    const promises = targets.map((target: AutotaskQuery) => {
      return this.queryData(target);
    });

    return Promise.all(promises);
  }

  async queryData(target: AutotaskQuery) {
    const response = await this.postResource('query', {
      queryType: target.queryType,
      filter: target.filter,
    });

    return {
      data: response,
      refId: target.refId,
    };
  }

  async testDatasource() {
    try {
      await this.postResource('test', {});
      return {
        status: 'success',
        message: 'Success: Connected to Autotask API',
      };
    } catch (error) {
      return {
        status: 'error',
        message: `Error: ${error.message}`,
      };
    }
  }
}
