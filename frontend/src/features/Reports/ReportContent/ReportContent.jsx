import dayjs from 'dayjs';
import { useMemo } from 'react';
import relativeTime from 'dayjs/plugin/relativeTime';
import { useMaterialReactTable } from 'material-react-table';

import RowHeader from '@common/RowHeader';
import { EmptyComponent } from '@common/utils';
import DataTable from '@common/DataTable/DataTable';
import { ASSETS_IN_REPORTS_HEADER } from '@features/Reports/constants';
import { Stack } from '@mui/material';

dayjs.extend(relativeTime);

export default function ReportContent({ sinceValue, assets }) {
  const columns = useMemo(() => ASSETS_IN_REPORTS_HEADER, []);

  const table = useMaterialReactTable({
    columns,
    data: assets,
    enableSorting: false,
    enableFilters: false,
    enableFullScreenToggle: false,
    enableHiding: false,
    enableTopToolbar: false,
    muiTableHeadRowProps: {
      sx: { padding: '1rem' },
    },
    muiTableBodyRowProps: {
      sx: { padding: '0.2rem' },
    },
    enableColumnResizing: true,
    enableDensityToggle: false,
    initialState: { density: 'compact' },
    renderEmptyRowsFallback: () => <EmptyComponent padding="1rem 1rem" />,
  });

  return (
    <Stack data-tour={'reports-7'}>
      <RowHeader title="Asset Details" caption={`Asset movement since ${dayjs(sinceValue).fromNow()}`} />
      <DataTable table={table} />
    </Stack>
  );
}
