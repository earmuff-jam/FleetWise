import { useEffect, useState } from 'react';

import { enqueueSnackbar } from 'notistack';
import { useParams } from 'react-router-dom';
import { useDispatch, useSelector } from 'react-redux';

import SimpleModal from '@common/SimpleModal';
import { Skeleton, Stack } from '@mui/material';
import { AddRounded } from '@mui/icons-material';

import { ConfirmationBoxModal } from '@common/utils';
import ItemHeader from '@common/ItemCard/ItemHeader/ItemHeader';
import { inventoryActions } from '@features/Assets/inventorySlice';
import { maintenancePlanItemActions } from '@features/MaintenancePlanItemDetails/maintenancePlanItemSlice';

import AddItem from '@common/ItemCard/AddItem/AddItem';
import ItemContent from '@common/ItemCard/ItemContent/ItemContent';
import ItemGraphWrapper from '@common/ItemCard/ItemGraph/ItemGraphWrapper';

export default function MaintenancePlanItemDetails() {
  const { id } = useParams();
  const dispatch = useDispatch();

  const {
    selectedMaintenancePlan,
    itemsInMaintenancePlan = [],
    selectedMaintenancePlanImage,
    loading = false,
  } = useSelector((state) => state.maintenancePlanItem);

  const [selectedIDList, setSelectedIDList] = useState([]);
  const [displayModal, setDisplayModal] = useState(false);
  const [openConfirmationBoxModal, setOpenConfirmationBoxModal] = useState(false);

  const handleOpenModal = () => {
    setDisplayModal(true);
    dispatch(inventoryActions.getAllInventoriesForUser());
  };

  const handleOpenConfirmationBoxModal = () => setOpenConfirmationBoxModal(!openConfirmationBoxModal);

  const resetConfirmationBoxModal = () => setOpenConfirmationBoxModal(false);

  const confirmDelete = () => {
    dispatch(
      maintenancePlanItemActions.removeItemsFromMaintenancePlan({
        id: selectedMaintenancePlan?.id,
        selectedIDList,
      })
    );
    enqueueSnackbar(`Removed association of assets for ${selectedMaintenancePlan.name}.`, {
      variant: 'default',
    });
    setSelectedIDList([]);
    resetConfirmationBoxModal();
  };

  const addItems = () => {
    const collaborators = selectedMaintenancePlan.sharable_groups;
    dispatch(
      maintenancePlanItemActions.addItemsInPlan({
        id: selectedMaintenancePlan?.id,
        selectedIDList,
        collaborators,
      })
    );
    enqueueSnackbar(`Added association of assets for ${selectedMaintenancePlan.name}.`, {
      variant: 'success',
    });
    resetSelection();
  };

  const handleRemoveAssociation = () => {
    handleOpenConfirmationBoxModal();
  };

  const resetSelection = () => {
    setDisplayModal(false);
    setSelectedIDList([]);
  };

  useEffect(() => {
    if (!loading) {
      dispatch(maintenancePlanItemActions.getSelectedImage({ id }));
    }
  }, [id, loading]);

  useEffect(() => {
    if (id) {
      dispatch(maintenancePlanItemActions.getItemsInMaintenancePlan(id));
      dispatch(maintenancePlanItemActions.getSelectedMaintenancePlan(id));
    }
  }, [id]);

  if (loading) {
    return <Skeleton height="20rem" />;
  }

  return (
    <Stack spacing={3} data-tour="selected-plan-0">
      <ItemHeader
        label={selectedMaintenancePlan?.name ? `${selectedMaintenancePlan.name} Overview` : 'Maintenance Plan Overview'}
        caption="View details of selected maintenance plan"
        item={selectedMaintenancePlan}
        image={selectedMaintenancePlanImage}
        favBtnDataTour="selected-plan-1"
        shareBtnDataTour="selected-plan-2"
        imageBtnDataTour="selected-plan-3"
      />
      <ItemContent
        selectedIDList={selectedIDList}
        setSelectedIDList={setSelectedIDList}
        items={itemsInMaintenancePlan}
        handleOpenModal={handleOpenModal}
        handleRemoveAssociation={handleRemoveAssociation}
        tableDataTour="selected-plan-6"
        primaryBtnDataTour="selected-plan-4"
        secondaryBtnDataTour="selected-plan-5"
      />
      <ItemGraphWrapper associatedAssets={itemsInMaintenancePlan} graphDataTour="selected-plan-7" />
      {displayModal && (
        <SimpleModal
          title={`Add items to ${selectedMaintenancePlan?.name}`}
          handleClose={resetSelection}
          maxSize="md"
          showSecondaryButton
          secondaryButtonAction={addItems}
          secondaryButtonIcon={<AddRounded />}
          disableSecondaryButton={selectedIDList.length <= 0}
        >
          <AddItem
            selectedIDList={selectedIDList}
            setSelectedIDList={setSelectedIDList}
            resetSelection={resetSelection}
            associatedItems={itemsInMaintenancePlan}
          />
        </SimpleModal>
      )}
      <ConfirmationBoxModal
        openDialog={openConfirmationBoxModal}
        title="Confirm deletion"
        handleClose={resetConfirmationBoxModal}
        maxSize="xs"
        confirmDelete={confirmDelete}
      />
    </Stack>
  );
}
