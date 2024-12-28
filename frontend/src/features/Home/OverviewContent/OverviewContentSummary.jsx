import { Stack } from '@mui/material';
import { prefix } from '@common/utils';

import { Inventory2Rounded, PaidRounded } from '@mui/icons-material';

import RowHeader from '@common/RowHeader';
import OverviewCardItem from '@features/Home/OverviewCard/OverviewCardItem';
import OverviewCardWrapper from '@features/Home/OverviewCard/OverviewCardWrapper';

export default function OverviewContentSummary({ assets = [] }) {
  const zeroCostItems = assets?.filter((v) => v.price === 0);
  const totalAssetCosts = assets?.reduce((acc, el) => {
    acc += el?.price || 0;
    return acc;
  }, 0);

  return (
    <Stack data-tour="overview-2">
      <RowHeader title="Cost Summary" caption="View details about associated costs and unestimated assets." />
      <Stack direction="row" spacing={2} useFlexGap flexWrap="wrap">
        <OverviewCardWrapper>
          <OverviewCardItem
            label="Total estimated cost"
            icon={<PaidRounded />}
            color="success.main"
            dataLabel={prefix('$', totalAssetCosts)}
            word=""
            tooltipTitle={'Total approximate dollar value of all assets.'}
          />
        </OverviewCardWrapper>
        <OverviewCardWrapper>
          <OverviewCardItem
            label="Total Unestimated assets"
            icon={<Inventory2Rounded />}
            color="text.secondary"
            dataLabel={`${zeroCostItems?.length || 0}`}
            word="asset"
            tooltipTitle={'Total assets that do not have cost price associated with them.'}
          />
        </OverviewCardWrapper>
      </Stack>
    </Stack>
  );
}
