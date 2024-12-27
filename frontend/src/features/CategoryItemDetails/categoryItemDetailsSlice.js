import { createSlice } from '@reduxjs/toolkit';

const initialState = {
  error: '',
  loading: false,
  selectedCategory: {},
  selectedCategoryImage: '',
  itemsInCategory: [],
};

const categoryItemDetailsSlice = createSlice({
  name: 'categoryItemDetails',
  initialState,
  reducers: {
    getCategory: (state) => {
      state.loading = true;
      state.error = '';
      state.selectedCategory = {};
    },
    getCategorySuccess: (state, action) => {
      state.selectedCategory = action.payload;
      state.loading = false;
      state.error = '';
    },
    getCategoryFailure: (state) => {
      state.loading = false;
      state.error = '';
      state.selectedCategory = {};
    },
    updateCollaborators: (state) => {
      state.loading = true;
      state.error = '';
    },
    updateCollaboratorsSuccess: (state, action) => {
      state.selectedCategory = action.payload;
      state.loading = false;
      state.error = '';
    },
    updateCollaboratorsFailure: (state) => {
      state.loading = false;
      state.error = '';
    },
    getItemsForCategory: (state) => {
      state.loading = true;
      state.error = '';
      state.itemsInCategory = [];
    },
    getItemsForCategorySuccess: (state, action) => {
      state.itemsInCategory = action.payload;
      state.loading = false;
      state.error = '';
    },
    getItemsForCategoryFailure: (state) => {
      state.loading = false;
      state.error = '';
      state.itemsInCategory = [];
    },
    addItemsInCategory: (state) => {
      state.loading = true;
      state.error = '';
    },
    addItemsInCategorySuccess: (state, action) => {
      state.itemsInCategory = [...action.payload];
      state.loading = false;
      state.error = '';
    },
    addItemsInCategoryFailure: (state) => {
      state.loading = false;
      state.error = '';
      state.itemsInCategory = [];
    },
    removeItemsFromCategory: (state) => {
      state.loading = true;
      state.error = '';
    },
    removeItemsFromCategorySuccess: (state, action) => {
      const removedItems = action.payload;
      const filteredItemsFromCategoryList = state.itemsInCategory.filter((v) => !removedItems.includes(v.id));
      state.loading = false;
      state.error = '';
      state.itemsInCategory = [...filteredItemsFromCategoryList];
    },
    removeItemsFromCategoryFailure: (state) => {
      state.loading = false;
      state.error = '';
    },
    uploadImage: (state) => {
      state.error = '';
      state.loading = true;
    },
    uploadImageSuccess: (state) => {
      state.loading = false;
      state.selectedCategoryImage = '';
      state.error = '';
    },
    uploadImageFailure: (state) => {
      state.loading = false;
      state.error = '';
    },
    getSelectedImage: (state) => {
      state.error = '';
    },
    getSelectedImageSuccess: (state, action) => {
      state.selectedCategoryImage = action.payload;
      state.loading = false;
      state.error = '';
    },
    getSelectedImageFailure: (state) => {
      state.loading = false;
      state.error = '';
      state.selectedCategoryImage = '';
    },
  },
});

export const { actions: categoryItemDetailsActions } = categoryItemDetailsSlice;
export default categoryItemDetailsSlice.reducer;
