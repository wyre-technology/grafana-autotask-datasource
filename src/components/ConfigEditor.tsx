import React from 'react';
import { DataSourcePluginOptionsEditorProps } from '@grafana/data';
import { InlineField, Input, SecretInput } from '@grafana/ui';
import { AutotaskDataSourceOptions } from '../types';

interface Props extends DataSourcePluginOptionsEditorProps<AutotaskDataSourceOptions> {}

export function ConfigEditor(props: Props) {
  const { options, onOptionsChange } = props;

  const onUsernameChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    onOptionsChange({
      ...options,
      jsonData: {
        ...options.jsonData,
        username: event.target.value,
      },
    });
  };

  const onSecretChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    onOptionsChange({
      ...options,
      secureJsonData: {
        ...options.secureJsonData,
        secret: event.target.value,
      },
    });
  };

  const onIntegrationCodeChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    onOptionsChange({
      ...options,
      secureJsonData: {
        ...options.secureJsonData,
        integrationCode: event.target.value,
      },
    });
  };

  const onZoneChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    onOptionsChange({
      ...options,
      jsonData: {
        ...options.jsonData,
        zone: event.target.value,
      },
    });
  };

  return (
    <div className="gf-form-group">
      <div className="gf-form">
        <InlineField label="Username" labelWidth={12}>
          <Input
            value={options.jsonData.username || ''}
            onChange={onUsernameChange}
            placeholder="Enter your Autotask username"
            width={40}
          />
        </InlineField>
      </div>
      <div className="gf-form">
        <InlineField label="Secret" labelWidth={12}>
          <SecretInput
            value={options.secureJsonData?.secret || ''}
            onChange={onSecretChange}
            placeholder="Enter your Autotask secret"
            width={40}
          />
        </InlineField>
      </div>
      <div className="gf-form">
        <InlineField label="Integration Code" labelWidth={12}>
          <SecretInput
            value={options.secureJsonData?.integrationCode || ''}
            onChange={onIntegrationCodeChange}
            placeholder="Enter your Autotask integration code"
            width={40}
          />
        </InlineField>
      </div>
      <div className="gf-form">
        <InlineField label="Zone" labelWidth={12}>
          <Input
            value={options.jsonData.zone || ''}
            onChange={onZoneChange}
            placeholder="Enter your Autotask zone (e.g., webservices14.autotask.net)"
            width={40}
          />
        </InlineField>
      </div>
    </div>
  );
}
