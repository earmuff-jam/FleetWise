import { Stack } from '@mui/material';

import { CategoryRounded, EngineeringRounded, WarningRounded } from '@mui/icons-material';

import RowHeader from '@common/RowHeader';
import OverviewCardItem from '@features/Home/OverviewCard/OverviewCardItem';
import OverviewCardWrapper from '@features/Home/OverviewCard/OverviewCardWrapper';

export default function OverviewHeader({
  assetsUnderCategories = [],
  assetsUnderMaintenancePlans = [],
  assetsPastDue = [],
}) {
  return (
    <Stack data-tour="overview-1">
      <RowHeader title="Asset Summary" caption="View details about your assets." />
      <Stack direction="row" spacing={2} useFlexGap flexWrap="wrap">
        <OverviewCardWrapper>
          <OverviewCardItem
            label="assigned categories"
            icon={<CategoryRounded />}
            color="info.main"
            tooltipTitle={assetsUnderCategories.map((v) => v.name).join(', ')}
            dataLabel={assetsUnderCategories.flatMap((v) => v.items).length}
          />
        </OverviewCardWrapper>
        <OverviewCardWrapper>
          <OverviewCardItem
            label="assigned maintenance plan"
            icon={<EngineeringRounded />}
            color="info.main"
            tooltipTitle={assetsUnderMaintenancePlans.map((v) => v.name).join(', ')}
            dataLabel={assetsUnderMaintenancePlans.flatMap((v) => v.items).length}
          />
        </OverviewCardWrapper>
        <OverviewCardWrapper>
          <OverviewCardItem
            label="past return deadline"
            icon={<WarningRounded />}
            color="error.main"
            tooltipTitle={assetsPastDue.join(', ')}
            dataLabel={assetsPastDue.length}
          />
        </OverviewCardWrapper>
      </Stack>
    </Stack>
  );
}
