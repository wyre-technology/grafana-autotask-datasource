import { DataSourceJsonData } from '@grafana/data';
import { DataQuery } from '@grafana/schema';

export type AutotaskEntityType = 'tickets' | 'resources' | 'companies' | 'contacts';

export interface AutotaskQuery extends DataQuery {
  queryType: AutotaskEntityType;
  filter: string;
  timeField: string;
  maxRecords: number;
}

export const DEFAULT_AUTOTASK_QUERY: Partial<AutotaskQuery> = {
  queryType: 'tickets',
  filter: '',
  timeField: '',
  maxRecords: 500,
};

export interface AutotaskDatasourceOptions extends DataSourceJsonData {
  username: string;
  url: string;
}

export interface AutotaskSecureJsonData {
  // Optional: Grafana secure-json fields are write-only — never read back
  // from the server, absent on config load (only secureJsonFields booleans
  // indicate whether a value is set). The ConfigEditor builds partial
  // updates ({ ...secureJsonData, secret }) so requiring both fields made
  // TS reject the partial as not-assignable (TS2322). Readers already use
  // secureJsonData?.secret ?? '' (optional-safe).
  secret?: string;
  integrationCode?: string;
}

// Entity type metadata for the query editor
export const ENTITY_TYPES: Array<{
  label: string;
  value: AutotaskEntityType;
  description: string;
  timeFields: string[];
}> = [
  {
    label: 'Tickets',
    value: 'tickets',
    description: 'Service desk tickets',
    timeFields: ['createDate', 'dueDateTime', 'lastActivityDate', 'completedDate'],
  },
  {
    label: 'Companies',
    value: 'companies',
    description: 'Customer companies',
    timeFields: ['createDate', 'lastActivityDate'],
  },
  {
    label: 'Contacts',
    value: 'contacts',
    description: 'Company contacts',
    timeFields: ['createDate', 'lastActivityDate'],
  },
  {
    label: 'Resources',
    value: 'resources',
    description: 'Internal resources/technicians',
    timeFields: [],
  },
];
