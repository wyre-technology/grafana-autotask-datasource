import React from 'react';
import { Select, InlineField, Input } from '@grafana/ui';
import { QueryEditorProps, SelectableValue } from '@grafana/data';
import { AutotaskDatasource } from '../datasource';
import {
  AutotaskQuery,
  AutotaskDatasourceOptions,
  AutotaskEntityType,
  ENTITY_TYPES,
  DEFAULT_AUTOTASK_QUERY,
} from '../types';

type Props = QueryEditorProps<AutotaskDatasource, AutotaskQuery, AutotaskDatasourceOptions>;

export function QueryEditor({ query, onChange, onRunQuery }: Props) {
  const q: AutotaskQuery = {
    ...DEFAULT_AUTOTASK_QUERY,
    ...query,
  } as AutotaskQuery;

  const entityMeta = ENTITY_TYPES.find((e) => e.value === q.queryType);

  const entityOptions: Array<SelectableValue<AutotaskEntityType>> = ENTITY_TYPES.map((e) => ({
    label: e.label,
    value: e.value,
    description: e.description,
  }));

  const timeFieldOptions: Array<SelectableValue<string>> = [
    { label: 'None', value: '', description: 'Do not filter by time range' },
    ...(entityMeta?.timeFields.map((f) => ({ label: f, value: f })) || []),
  ];

  const onQueryTypeChange = (value: SelectableValue<AutotaskEntityType>) => {
    if (value.value) {
      onChange({ ...q, queryType: value.value, timeField: '' });
      onRunQuery();
    }
  };

  const onFilterChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    onChange({ ...q, filter: event.target.value });
  };

  const onFilterBlur = () => {
    onRunQuery();
  };

  const onTimeFieldChange = (value: SelectableValue<string>) => {
    onChange({ ...q, timeField: value.value || '' });
    onRunQuery();
  };

  return (
    <div>
      <div className="gf-form-inline">
        <InlineField label="Entity" labelWidth={12} tooltip="The Autotask entity type to query">
          <Select
            options={entityOptions}
            value={entityOptions.find((o) => o.value === q.queryType)}
            onChange={onQueryTypeChange}
            width={24}
          />
        </InlineField>
        <InlineField
          label="Time Field"
          labelWidth={12}
          tooltip="Map a date field to the Grafana time range for filtering"
        >
          <Select
            options={timeFieldOptions}
            value={timeFieldOptions.find((o) => o.value === q.timeField)}
            onChange={onTimeFieldChange}
            width={24}
            isClearable
          />
        </InlineField>
      </div>
      <div className="gf-form-inline">
        <InlineField
          label="Filter"
          labelWidth={12}
          tooltip="Autotask query filter JSON (optional). Example: {&quot;op&quot;:&quot;eq&quot;,&quot;field&quot;:&quot;status&quot;,&quot;value&quot;:1}"
          grow
        >
          <Input
            value={q.filter}
            placeholder='{"op":"eq","field":"status","value":1}'
            onChange={onFilterChange}
            onBlur={onFilterBlur}
          />
        </InlineField>
      </div>
    </div>
  );
}
