import React from 'react';
import { DataSourcePluginOptionsEditorProps } from '@grafana/data';
import { AutotaskDatasourceOptions, AutotaskSecureJsonData } from '../types';
import { InlineField, Input, SecretInput, FieldSet } from '@grafana/ui';

interface Props extends DataSourcePluginOptionsEditorProps<AutotaskDatasourceOptions, AutotaskSecureJsonData> {}

export function ConfigEditor(props: Props) {
  const { options, onOptionsChange } = props;
  const { jsonData, secureJsonData, secureJsonFields } = options;

  const onURLChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    let url = event.target.value.trim();

    if (url) {
      url = url.replace(/^http:\/\//, 'https://');
      if (!url.startsWith('https://')) {
        url = 'https://' + url;
      }
      url = url.replace(/\/+$/, '');
    }

    onOptionsChange({
      ...options,
      jsonData: { ...jsonData, url },
    });
  };

  const onUsernameChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    onOptionsChange({
      ...options,
      jsonData: { ...jsonData, username: event.target.value },
    });
  };

  const onSecretChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    onOptionsChange({
      ...options,
      secureJsonData: {
        ...secureJsonData,
        secret: event.target.value,
      },
    });
  };

  const onIntegrationCodeChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    onOptionsChange({
      ...options,
      secureJsonData: {
        ...secureJsonData,
        integrationCode: event.target.value,
      },
    });
  };

  return (
    <>
      <FieldSet label="Connection">
        <InlineField
          label="API URL"
          labelWidth={14}
          tooltip="The Autotask API base URL, e.g. https://webservices6.autotask.net. Leave as default if unsure — the zone is auto-detected."
        >
          <Input
            value={jsonData.url || ''}
            placeholder="https://webservices.autotask.net"
            onChange={onURLChange}
            width={40}
          />
        </InlineField>
        <InlineField label="Username" labelWidth={14} tooltip="Your Autotask API username (email address)">
          <Input
            value={jsonData.username || ''}
            placeholder="user@company.com"
            onChange={onUsernameChange}
            width={40}
          />
        </InlineField>
      </FieldSet>

      <FieldSet label="Authentication">
        <InlineField label="API Secret" labelWidth={14} tooltip="Your Autotask API secret key">
          <SecretInput
            isConfigured={!!secureJsonFields?.secret}
            value={secureJsonData?.secret || ''}
            placeholder="API secret"
            width={40}
            onReset={() => {
              onOptionsChange({
                ...options,
                secureJsonFields: { ...secureJsonFields, secret: false },
                secureJsonData: { ...secureJsonData, secret: '' },
              });
            }}
            onChange={onSecretChange}
          />
        </InlineField>
        <InlineField
          label="Integration Code"
          labelWidth={14}
          tooltip="Your Autotask API integration code"
        >
          <SecretInput
            isConfigured={!!secureJsonFields?.integrationCode}
            value={secureJsonData?.integrationCode || ''}
            placeholder="Integration code"
            width={40}
            onReset={() => {
              onOptionsChange({
                ...options,
                secureJsonFields: { ...secureJsonFields, integrationCode: false },
                secureJsonData: { ...secureJsonData, integrationCode: '' },
              });
            }}
            onChange={onIntegrationCodeChange}
          />
        </InlineField>
      </FieldSet>
    </>
  );
}
