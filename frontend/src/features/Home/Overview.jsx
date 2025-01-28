import dayjs from 'dayjs';

import { useEffect } from 'react';

import { useDispatch, useSelector } from 'react-redux';

import { Skeleton, Stack } from '@mui/material';

import RowHeader from '@common/RowHeader';
import { assetSummaryActions } from '@features/Home/SummarySlice';
import OverviewHeader from '@features/Home/OverviewHeader/OverviewHeader';
import OverviewContent from '@features/Home/OverviewContent/OverviewContent';

export default function Overview() {
  const dispatch = useDispatch();
  const { assetSummary = [], loading } = useSelector((state) => state.summary);

  const fetchAssignedAssets = (assets = [], filterParam) => {
    // assets within filterParam must be present
    assets?.filter((asset) => asset.type.toUpperCase() === filterParam && asset.items[0] != '');
  };

  const assetsPastDue = assetSummary?.AssetList?.reduce((acc, el) => {
    if (dayjs(el.returntime).isBefore(dayjs())) {
      acc.push(el.name);
    }
    return acc;
  }, []);

  useEffect(() => {
    dispatch(assetSummaryActions.getAssetSummary());
  }, []);

  if (loading) {
    return <Skeleton height="50vh" />;
  }

  return (
    <Stack data-tour="overview-0">
      <RowHeader title="Overview" caption="View a summarized report about your assets." />
      <Stack spacing={2}>
        <OverviewHeader
          assetsUnderCategories={fetchAssignedAssets(assetSummary?.AssetSummaryList, 'C')}
          assetsUnderMaintenancePlans={fetchAssignedAssets(assetSummary?.AssetSummaryList, 'M')}
          assetsPastDue={assetsPastDue}
        />
        <OverviewContent
          assets={assetSummary?.AssetList || []}
          categories={assetSummary?.AssetSummaryList?.filter((v) => v.type.toUpperCase() === 'C')}
          maintenancePlans={assetSummary?.AssetSummaryList?.filter((v) => v.type.toUpperCase() === 'M')}
        />
      </Stack>
    </Stack>
  );
}
