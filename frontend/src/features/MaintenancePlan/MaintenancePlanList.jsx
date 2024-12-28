import { useEffect, useState } from 'react';

import { useDispatch, useSelector } from 'react-redux';

import { Box } from '@mui/material';
import SimpleModal from '@common/SimpleModal';
import SectionCardHeader from '@common/SectionCard/SectionCardHeader';
import SectionCardContent from '@common/SectionCard/SectionCardContent';
import AddMaintenancePlan from '@features/MaintenancePlan/AddMaintenancePlan';
import { maintenancePlanActions } from '@features/MaintenancePlan/maintenanceSlice';

const MaintenancePlanList = () => {
  const dispatch = useDispatch();
  const { maintenancePlan, loading } = useSelector((state) => state.maintenance);

  const [sortedData, setSortedData] = useState([]);
  const [displayModal, setDisplayModal] = useState(false);

  const [selectedFilter, setSelectedFilter] = useState('');
  const [sortingOrder, setSortingOrder] = useState(true); // false ascending
  const [selectedMaintenancePlanID, setSelectedMaintenancePlanID] = useState('');

  const toggleModal = () => setDisplayModal(!displayModal);
  const handleCloseAddNewPlan = () => {
    setDisplayModal(false);
    setSelectedMaintenancePlanID('');
  };

  const handleDownload = () => {
    dispatch(maintenancePlanActions.download());
  };

  const filterAndBuildMaintenancePlans = (plans, selectedFilter) => {
    if (selectedFilter.length > 0) {
      return plans.filter((element) => element.maintenance_status_name === selectedFilter);
    } else {
      return sortedData;
    }
  };

  const removeSelectedMaintenancePlan = (id) => dispatch(maintenancePlanActions.removePlan(id));

  useEffect(() => {
    if (sortingOrder) {
      const draft = [...maintenancePlan].sort((a, b) => new Date(a.updated_at) - new Date(b.updated_at));
      setSortedData(draft);
    } else {
      setSortedData(maintenancePlan);
    }
  }, [sortingOrder, maintenancePlan]);

  useEffect(() => {
    dispatch(maintenancePlanActions.getPlans(100));
  }, []);

  return (
    <Box sx={{ py: 2 }} data-tour="plans-0">
      <SectionCardHeader
        title="Maintenance Plans"
        caption={selectedFilter ? `Applying ${selectedFilter} status filter` : 'Assign items to maintenance plans'}
        primaryBtnTitle="Add plan"
        toggleModal={toggleModal}
        selectedFilter={selectedFilter}
        setSelectedFilter={setSelectedFilter}
        sortingOrder={sortingOrder}
        setSortingOrder={setSortingOrder}
        handleDownload={handleDownload}
        addBtnDataTour={'plans-1'}
        downloadBtnDataTour={'plans-2'}
        filterBtnDataTour={'plans-3'}
        sortBtnDataTour={'plans-4'}
        disableDownloadIcon={Boolean(maintenancePlan) && maintenancePlan.length <= 0}
      />
      <SectionCardContent
        loading={loading}
        prefixURI={'plan'}
        displayModal={displayModal}
        setDisplayModal={setDisplayModal}
        setSelectedID={setSelectedMaintenancePlanID}
        removeItem={removeSelectedMaintenancePlan}
        content={filterAndBuildMaintenancePlans(maintenancePlan, selectedFilter)}
      />
      {displayModal && (
        <SimpleModal
          title="Add new maintenance plan"
          subtitle="Create maintenance plan to associate assets and periodically perform checks on them.
"
          handleClose={handleCloseAddNewPlan}
          maxSize="sm"
        >
          <AddMaintenancePlan
            maintenancePlan={maintenancePlan}
            handleCloseAddNewPlan={handleCloseAddNewPlan}
            selectedMaintenancePlanID={selectedMaintenancePlanID}
            setSelectedMaintenancePlanID={setSelectedMaintenancePlanID}
          />
        </SimpleModal>
      )}
    </Box>
  );
};

export default MaintenancePlanList;
