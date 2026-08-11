import { test, expect } from '@grafana/plugin-e2e';

const FILTER_PLACEHOLDER = '{"op":"eq","field":"status","value":1}';

test('smoke: should render query editor', async ({ panelEditPage, readProvisionedDataSource }) => {
  const ds = await readProvisionedDataSource({ fileName: 'datasources.yml' });
  await panelEditPage.datasource.set(ds.name);

  const queryEditorRow = panelEditPage.getQueryEditorRow('A');
  await expect(queryEditorRow.getByText('Entity', { exact: true })).toBeVisible();
  await expect(queryEditorRow.getByText('Time Field', { exact: true })).toBeVisible();
  await expect(queryEditorRow.getByPlaceholder(FILTER_PLACEHOLDER)).toBeVisible();
});

test('should trigger new query when the filter is changed', async ({ panelEditPage, readProvisionedDataSource }) => {
  const ds = await readProvisionedDataSource({ fileName: 'datasources.yml' });
  await panelEditPage.datasource.set(ds.name);

  const filter = panelEditPage.getQueryEditorRow('A').getByPlaceholder(FILTER_PLACEHOLDER);
  await filter.fill(FILTER_PLACEHOLDER);

  // The editor only runs the query on blur, so wire up the listener first.
  const queryReq = panelEditPage.waitForQueryDataRequest();
  await filter.blur();
  await expect(await queryReq).toBeTruthy();
});
