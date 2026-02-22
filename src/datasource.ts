import {
  DataSourceInstanceSettings,
  DataQueryRequest,
  DataQueryResponse,
} from '@grafana/data';
import { DataSourceWithBackend, getBackendSrv } from '@grafana/runtime';
import { AutotaskQuery, AutotaskDatasourceOptions } from './types';
import { lastValueFrom, from } from 'rxjs';
import { catchError } from 'rxjs/operators';

export class AutotaskDatasource extends DataSourceWithBackend<AutotaskQuery, AutotaskDatasourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<AutotaskDatasourceOptions>) {
    super(instanceSettings);
  }

  // DataSourceWithBackend handles query() automatically by proxying to the Go backend.
  // We only need to override if we want custom frontend logic.

  filterQuery(query: AutotaskQuery): boolean {
    return !!query.queryType;
  }

  async testDatasource() {
    // Use the resource endpoint for testing since it returns zone info
    try {
      const response = await lastValueFrom(
        from(
          getBackendSrv().post(
            `/api/datasources/uid/${this.uid}/resources/zoneinfo`
          )
        ).pipe(
          catchError((error) => {
            throw {
              status: 'error',
              message: `Failed to connect: ${error.statusText || error.message || 'Unknown error'}`,
            };
          })
        )
      );

      if (response?.zoneName || response?.ZoneName) {
        return {
          status: 'success',
          message: `Connected to Autotask (Zone: ${response.zoneName || response.ZoneName})`,
        };
      }

      throw new Error('No zone information returned');
    } catch (error: any) {
      return {
        status: 'error',
        message: error?.message || 'Unknown error connecting to Autotask',
      };
    }
  }
}
