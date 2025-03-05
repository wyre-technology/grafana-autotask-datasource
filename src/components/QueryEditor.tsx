import React from 'react';
import { QueryEditorProps } from '@grafana/data';
import { AutotaskQuery } from '../types';
import { Select, Input } from '@grafana/ui';

export function QueryEditor({ query, onChange, datasource }: QueryEditorProps<AutotaskQuery>) {
  const onQueryTypeChange = (value: string) => {
    onChange({ ...query, queryType: value });
  };

  const onFilterChange = (value: string) => {
    onChange({ ...query, filter: value });
  };

  return (
    <div className="gf-form">
      <div className="gf-form-inline">
        <div className="gf-form">
          <label className="gf-form-label width-8">Query Type</label>
          <Select
            value={query.queryType}
            onChange={(value) => onQueryTypeChange(value.value)}
            options={[
              { label: 'Tickets', value: 'tickets' },
              { label: 'Resources', value: 'resources' },
              { label: 'Companies', value: 'companies' },
            ]}
          />
        </div>
        <div className="gf-form">
          <label className="gf-form-label width-8">Filter</label>
          <Input
            value={query.filter}
            onChange={(e) => onFilterChange(e.currentTarget.value)}
            placeholder="Enter filter query..."
          />
        </div>
      </div>
    </div>
  );
}
