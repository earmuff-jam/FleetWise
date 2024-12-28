import { useEffect, useState } from 'react';
import { enqueueSnackbar } from 'notistack';
import { useParams } from 'react-router-dom';
import { useDispatch, useSelector } from 'react-redux';

import { Skeleton, Stack } from '@mui/material';
import { AddRounded } from '@mui/icons-material';

import SimpleModal from '@common/SimpleModal';
import { ConfirmationBoxModal } from '@common/utils';
import AddItem from '@common/ItemCard/AddItem/AddItem';
import ItemHeader from '@common/ItemCard/ItemHeader/ItemHeader';

import ItemContent from '@common/ItemCard/ItemContent/ItemContent';
import { inventoryActions } from '@features/Assets/inventorySlice';
import ItemGraphWrapper from '@common/ItemCard/ItemGraph/ItemGraphWrapper';
import { categoryItemDetailsActions } from '@features/CategoryItemDetails/categoryItemDetailsSlice';

export default function CategoryItemDetails() {
  const { id } = useParams();
  const dispatch = useDispatch();

  const {
    selectedCategory,
    selectedCategoryImage,
    itemsInCategory = [],
    loading = false,
  } = useSelector((state) => state.categoryItemDetails);

  const [selectedIDList, setSelectedIDList] = useState([]);
  const [displayModal, setDisplayModal] = useState(false);
  const [openConfirmationBoxModal, setOpenConfirmationBoxModal] = useState(false);

  const handleOpenModal = () => {
    setDisplayModal(true);
    dispatch(inventoryActions.getAllInventoriesForUser());
  };

  const handleOpenConfirmationBoxModal = () => setOpenConfirmationBoxModal(!openConfirmationBoxModal);

  const resetSelection = () => {
    setDisplayModal(false);
    setSelectedIDList([]);
  };

  const resetConfirmationBoxModal = () => setOpenConfirmationBoxModal(false);

  const confirmDelete = () => {
    dispatch(categoryItemDetailsActions.removeItemsFromCategory({ id: selectedCategory?.id, selectedIDList }));
    enqueueSnackbar(`Removed association of assets for ${selectedCategory.name}.`, {
      variant: 'default',
    });
    setSelectedIDList([]);
    resetConfirmationBoxModal();
  };

  const addItems = () => {
    const collaborators = selectedCategory.sharable_groups;
    dispatch(
      categoryItemDetailsActions.addItemsInCategory({ id: selectedCategory?.id, selectedIDList, collaborators })
    );
    enqueueSnackbar(`Added association of assets for ${selectedCategory.name}.`, {
      variant: 'success',
    });
    resetSelection();
  };

  const handleRemoveAssociation = () => {
    handleOpenConfirmationBoxModal();
  };

  useEffect(() => {
    if (!loading) {
      dispatch(categoryItemDetailsActions.getSelectedImage({ id }));
    }
  }, [loading]);

  useEffect(() => {
    if (id) {
      dispatch(categoryItemDetailsActions.getCategory(id));
      dispatch(categoryItemDetailsActions.getItemsForCategory(id));
    }
  }, [id]);

  if (loading) {
    return <Skeleton height="20rem" />;
  }

  return (
    <Stack spacing={3} data-tour="selected-category-0">
      <ItemHeader
        categoryMode
        label={selectedCategory?.name ? `${selectedCategory.name} Overview` : 'Category Overview'}
        caption="View details of selected category"
        item={selectedCategory}
        image={selectedCategoryImage}
        favBtnDataTour="selected-category-1"
        shareBtnDataTour="selected-category-2"
        imageBtnDataTour="selected-category-3"
      />
      <ItemContent
        selectedIDList={selectedIDList}
        setSelectedIDList={setSelectedIDList}
        items={itemsInCategory}
        handleOpenModal={handleOpenModal}
        handleRemoveAssociation={handleRemoveAssociation}
        tableDataTour="selected-category-6"
        primaryBtnDataTour="selected-category-4"
        secondaryBtnDataTour="selected-category-5"
      />
      <ItemGraphWrapper associatedAssets={itemsInCategory} graphDataTour="selected-category-7" />
      {displayModal && (
        <SimpleModal
          title={`Add items to ${selectedCategory?.name}`}
          handleClose={resetSelection}
          showSecondaryButton
          secondaryButtonAction={addItems}
          disableSecondaryButton={selectedIDList.length <= 0}
          secondaryButtonIcon={<AddRounded />}
          maxSize="md"
        >
          <AddItem
            selectedIDList={selectedIDList}
            setSelectedIDList={setSelectedIDList}
            resetSelection={resetSelection}
            associatedItems={itemsInCategory}
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
