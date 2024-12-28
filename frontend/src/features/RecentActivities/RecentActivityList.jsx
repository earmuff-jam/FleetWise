import dayjs from 'dayjs';
import { useDispatch } from 'react-redux';

import RowHeader from '@common/RowHeader';
import { Divider, Stack } from '@mui/material';
import relativeTime from 'dayjs/plugin/relativeTime';
import { DownloadRounded } from '@mui/icons-material';
import { profileActions } from '@features/Profile/profileSlice';
import RecentActivityAccordion from '@features/RecentActivities/RecentActivityAccordion/RecentActivityAccordion';

dayjs.extend(relativeTime);

const RecentActivityList = () => {
  const dispatch = useDispatch();
  const downloadRecentActivities = () => {
    const last30Days = dayjs().subtract(30, 'days').format('YYYY-MM-DDTHH:mm:ssZ');
    dispatch(profileActions.downloadRecentActivities({ last30Days }));
  };

  return (
    <Stack data-tour="recentActivities-0">
      <RowHeader
        title="Recent Activities"
        caption="View all of your recent activities here. You can view all associated activities performed by yourself. Track changes for the last ten assets. Download all recent activities for the past 30 days."
        primaryButtonTextLabel="Export"
        primaryStartIcon={<DownloadRounded />}
        handleClickPrimaryButton={downloadRecentActivities}
        primaryBtnDataTour="recentActivities-1"
      />
      <Divider sx={{ marginTop: '1rem', marginBottom: '1rem' }} />
      <RecentActivityAccordion />
    </Stack>
  );
};

export default RecentActivityList;
