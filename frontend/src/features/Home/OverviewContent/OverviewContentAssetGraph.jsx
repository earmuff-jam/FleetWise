import { Card, CardContent, Stack } from '@mui/material';

import RowHeader from '@common/RowHeader';
import PieChart from '@common/Chart/PieChart';

export default function OverviewContentAssetGraph({ assets = [], categories = [], maintenancePlans = [] }) {
  const formattedData = [categories, maintenancePlans, assets].map((v, index) => ({
    label: ['Categories', 'Plans', 'Assets'][index],
    count: v?.length || 0,
    backgroundColor: ['rgb(255, 99, 132)', 'rgb(54, 162, 235)', 'rgb(211, 211, 211)'][index],
    borderColor: ['rgb(255, 99, 132)', 'rgb(54, 162, 235)', 'rgb(211, 211, 211)'][index],
  }));

  return (
    <Stack data-tour="overview-4">
      <RowHeader title="Graph" caption="View details for asset summary" />
      <Card>
        <CardContent>
          <Stack>
            <Stack direction="row" spacing={2} sx={{ flexDirection: { xs: 'column', md: 'row' } }}>
              <PieChart height="15rem" legendLabel="assets summary" data={formattedData} />
            </Stack>
          </Stack>
        </CardContent>
      </Card>
    </Stack>
  );
}
