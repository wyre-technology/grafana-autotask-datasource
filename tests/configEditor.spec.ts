import { test, expect } from '@grafana/plugin-e2e';

test('smoke: should render config editor', async ({ createDataSourceConfigPage, readProvisionedDataSource, page }) => {
  const ds = await readProvisionedDataSource({ fileName: 'datasources.yml' });
  await createDataSourceConfigPage({ type: ds.type });

  await expect(page.getByPlaceholder('https://webservices.autotask.net')).toBeVisible();
  await expect(page.getByPlaceholder('user@company.com')).toBeVisible();
  await expect(page.getByPlaceholder('API secret')).toBeVisible();
  await expect(page.getByPlaceholder('Integration code')).toBeVisible();
});

// CI has no Autotask credentials, so the only health outcome we can assert
// deterministically is the offline one: the backend rejects incomplete
// settings in config.Validate() before any API call is made.
test('"Save & test" should fail when credentials are missing', async ({
  createDataSourceConfigPage,
  readProvisionedDataSource,
  page,
}) => {
  const ds = await readProvisionedDataSource({ fileName: 'datasources.yml' });
  const configPage = await createDataSourceConfigPage({ type: ds.type });

  await page.getByPlaceholder('user@company.com').fill('user@example.com');

  // AutotaskDatasource overrides testDatasource() to call its own
  // `resources/zoneinfo` endpoint rather than Grafana's health API, so there is
  // no health response to wait for. Assert on the resulting alert instead.
  await configPage.saveAndTest({ skipWaitForResponse: true });

  await expect(configPage).toHaveAlert('error');
});
