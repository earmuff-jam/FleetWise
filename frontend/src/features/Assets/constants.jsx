import dayjs from 'dayjs';

import QrCodeGen from '@common/QrCodeGen';
import relativeTime from 'dayjs/plugin/relativeTime';
import { CheckRounded, CloseRounded } from '@mui/icons-material';

dayjs.extend(relativeTime);

/**
 * Modal State is the display state of modal between adding a single inventory, adding a bulk
 * inventory and removing inventories
 *
 */
export const MODAL_STATE = {
  NONE: 'none',
  ADD_ITEM: 'item',
  BULK_ITEM: 'bulk',
  MORE_DETAILS: 'more',
  ASSIGN_CATEGORY: 'assign_category',
  ASSIGN_MAINTENANCE_PLAN: 'assign_maintenance_plan',
};

/**
 * Bulk Asset Headers for template download.
 */
export const BULK_ASSETS_HEADERS = [
  { label: 'Name' },
  { label: 'Description' },
  { label: 'Price' },
  { label: 'Quantity' },
  { label: 'Storage Location' },
  { label: 'Color' },
  { label: 'SKU' },
  { label: 'Barcode' },
  { label: 'Purchase Location' },
  { label: 'Minimum Height' },
  { label: 'Maximum Height' },
  { label: 'Minimum Weight' },
  { label: 'Maximum Weight' },
];

/**
 * default inventories landing page items to encourage users to use
 * various features of the application
 */
export const DEFAULT_INVENTORIES_LANDING_PAGE_TEXT = [
  {
    name: 'Create and view maintenance request',
    description: 'Create periodic maintenance of inventory items',
    imageSrc: 'images/books.jpg',
    imageAlt: 'default-maintenance-inventory-items-img',
    href: '/inventories/maintenance/list',
  },
  {
    name: 'Track record of inventories',
    description: 'Visualize items which are due for maintenance or need attention',
    imageSrc: 'images/items.jpg',
    imageAlt: 'default-track-record-inventory-items-img',
    href: 'inventories/forecast/list',
  },
  {
    name: 'Categories inventories and stay in sync',
    description: 'Understand your inventories better. Keep tabs on more important inventories',
    imageSrc: 'images/kitchen-items.jpg',
    imageAlt: 'default-bookmarked-inventory-items-img',
    href: '/inventories/categories/list',
  },
];

/**
 * combines the configs based on the selectedRow parameter. Returns a
 * list of tableRows derieved from the selectedRow object. If modifier function is passed
 * in, we build the table in accordance to it as well.
 *
 * @param {Object} selectedRow - the current selected row to build the table for
 * @returns {Array} tableRows - the combined row with modifiers applied if passed in.
 */
export const BUILD_TABLE_CONSTANTS = (columnLabels) => (selectedRow) => {
  if (!selectedRow) {
    return [];
  }
  const tableRows = columnLabels.map(({ id, colName, label, modifier }) => {
    let value = selectedRow[colName];
    if (modifier) {
      if (colName === 'qr_code') {
        value = modifier(selectedRow['name']);
      } else {
        value = modifier(value, { colName, label });
      }
    } else {
      value = value || 'N/A';
    }
    return {
      id,
      colName,
      label,
      value,
    };
  });

  return tableRows;
};

/**
 * INVENTORY LIST HEADERS STATIC COMPONENT
 * displayConcise lets users view the column name in bookmarked inventories
 * modifier fn lets the value be modified, for eg date will be modified with this property
 */
export const VIEW_INVENTORY_LIST_HEADERS = {
  name: {
    id: 1,
    colName: 'name',
    label: 'Name',
    displayConcise: true,
    modifier: (value) => `${value || '-'}`,
  },
  description: {
    id: 2,
    colName: 'description',
    label: 'Description',
    displayName: 'Description',
    displayConcise: false,
    modifier: (value) => `${value || '-'}`,
  },
  price: {
    id: 3,
    colName: 'price',
    label: 'Cost',
    displayConcise: true,
    modifier: (value) => `${+value || '-'}`,
  },
  barcode: {
    id: 5,
    colName: 'barcode',
    label: 'Barcode',
    modifier: (value) => `${value || '-'}`,
  },
  sku: {
    id: 6,
    colName: 'sku',
    label: 'SKU',
    modifier: (value) => `${value || '-'}`,
  },
  color: {
    id: 7,
    colName: 'color',
    label: 'Asset Color',
    modifier: (value) => `${value || '-'}`,
  },
  quantity: {
    id: 8,
    colName: 'quantity',
    label: 'Qty',
    displayConcise: true,
    modifier: (value) => `${+value || '-'}`,
  },
  location: {
    id: 9,
    colName: 'location',
    label: 'Location',
    displayConcise: true,
    modifier: (value) => `${value || '-'}`,
  },
  is_returnable: {
    id: 10,
    colName: 'is_returnable',
    label: 'Returnable',
    modifier: (value) => (value ? <CheckRounded color="primary" /> : <CloseRounded color="warning" />),
  },
  return_location: {
    id: 11,
    colName: 'return_location',
    label: 'Return Location',
    modifier: (value) => `${value || '-'}`,
  },
  max_weight: {
    id: 12,
    colName: 'max_weight',
    label: 'Max Weight',
    modifier: (value) => `${value || '-'}`,
  },
  min_weight: {
    id: 13,
    colName: 'min_weight',
    label: 'Min Weight',
    modifier: (value) => `${value || '-'}`,
  },
  max_height: {
    id: 14,
    colName: 'max_height',
    label: 'Max Height',
    modifier: (value) => `${value || '-'}`,
  },
  min_height: {
    id: 15,
    colName: 'min_height',
    label: 'Min Height',
    modifier: (value) => `${value || '-'}`,
  },
  bought_at: {
    id: 16,
    colName: 'bought_at',
    label: 'Purchase Location',
    modifier: (value) => `${value || '-'}`,
  },
  qr_code: {
    id: 17,
    colName: 'qr_code',
    label: 'QR Code',
    modifier: (value) => <QrCodeGen value={value} />,
  },
  updated_at: {
    id: 18,
    colName: 'updated_at',
    label: 'Updated',
    displayConcise: true,
    modifier: (value) => `${dayjs(value).fromNow()}`,
  },
  updator_name: {
    id: 19,
    colName: 'updator',
    label: 'Updated By',
    displayConcise: true,
    modifier: (value) => `${value || '-'}`,
  },
};

// blank form to add inventory details
export const BLANK_INVENTORY_FORM = {
  name: {
    id: 'name',
    label: 'Asset name',
    name: 'name',
    size: 'small',
    placeholder: 'Asset name',
    value: '',
    type: 'text',
    isRequired: true,
    errorMsg: '',
    validators: [
      {
        validate: (value) => value.trim().length === 0,
        message: 'Item name is required',
      },
      {
        validate: (value) => value.trim().length >= 200,
        message: 'Item name should be less than 50 characters',
      },
    ],
  },
  description: {
    id: 'description',
    label: 'Asset description',
    value: '',
    name: 'description',
    placeholder: 'Asset description in less than 500 characters',
    type: 'text',
    size: 'small',
    isRequired: true,
    errorMsg: '',
    validators: [
      {
        validate: (value) => value.trim().length === 0,
        message: 'Item description is required',
      },
      {
        validate: (value) => value.trim().length >= 500,
        message: 'Item description should be less than 50 characters',
      },
    ],
  },
  price: {
    id: 'price',
    label: 'Asset price (per unit)',
    value: '',
    name: 'price',
    size: 'small',
    placeholder: 'Asset price per unit',
    type: 'number',
    isRequired: false,
    errorMsg: '',
    validators: [
      {
        validate: (value) => value.trim().length === 0,
        message: 'Price for the selected item is required',
      },
      {
        validate: (value) => isNaN(value) || parseInt(value) <= 0,
        message: 'A positive number is required',
      },
      {
        validate: (value) => parseInt(value) > Number.MAX_SAFE_INTEGER,
        message: 'Value cannot exceed limit',
      },
    ],
  },
  barcode: {
    id: 'barcode',
    name: 'barcode',
    label: 'Barcode',
    value: '',
    type: 'text',
    size: 'small',
    placeholder: 'Barcode of the selected asset',
    isRequired: false,
    errorMsg: '',
    validators: [
      {
        validate: (value) => value.trim().length >= 50,
        message: 'Item barcode should be less than 50 characters',
      },
    ],
  },
  sku: {
    id: 'sku',
    name: 'sku',
    label: 'SKU',
    value: '',
    type: 'text',
    size: 'small',
    placeholder: 'SKU of the selected asset',
    isRequired: false,
    errorMsg: '',
    validators: [
      {
        validate: (value) => value.trim().length >= 50,
        message: 'SKU of item should be less than 50 characters',
      },
    ],
  },
  quantity: {
    id: 'quantity',
    name: 'quantity',
    label: 'Asset quantity',
    value: '',
    type: 'number',
    placeholder: 'The quantity of the asset',
    size: 'small',
    isRequired: true,
    errorMsg: '',
    validators: [
      {
        validate: (value) => value.trim().length === 0,
        message: 'Quantity for the selected item is required',
      },
      {
        validate: (value) => isNaN(value) || parseInt(value) <= 0,
        message: 'A positive number is required',
      },
      {
        validate: (value) => parseInt(value) > Number.MAX_SAFE_INTEGER,
        message: 'Value cannot exceed limit',
      },
    ],
  },
  bought_at: {
    id: 'bought_at',
    name: 'bought_at',
    label: 'Place of purchase',
    placeholder: 'Where did you buy this item?',
    value: '',
    type: 'text',
    size: 'small',
    isRequired: false,
    errorMsg: '',
    validators: [],
  },
  location: {
    id: 'location',
    label: 'Storage location',
    placeholder: 'Where do you want to store the item',
    type: 'text',
    isRequired: false,
    errorMsg: '',
    validators: [],
  },
  is_bookmarked: {
    id: 'is_bookmarked',
    value: false,
    type: 'boolean',
    isRequired: false,
    errorMsg: '',
    validators: [],
  },
  is_returnable: {
    id: 'is_returnable',
    value: false,
    type: 'boolean',
    isRequired: false,
    errorMsg: '',
    validators: [],
  },
  return_location: {
    id: 'return_location',
    name: 'return_location',
    label: 'Return location',
    placeholder: 'Where to return the item',
    value: '',
    type: 'text',
    size: 'small',
    isRequired: false,
    errorMsg: '',
    validators: [
      {
        validate: (value) => value.trim().length >= 50,
        message: 'Return location should be less than 50 characters',
      },
    ],
  },
  max_weight: {
    id: 'max_weight',
    name: 'max_weight',
    label: 'Maxmimum weight',
    placeholder: 'Maximum weight of the selected asset in kg',
    value: '',
    type: 'number',
    size: 'small',
    isRequired: false,
    errorMsg: '',
    validators: [
      {
        validate: (value) => isNaN(value) || parseInt(value) <= 0,
        message: 'A positive number is required',
      },
      {
        validate: (value) => parseInt(value) > Number.MAX_SAFE_INTEGER,
        message: 'Value cannot exceed limit',
      },
    ],
  },
  min_weight: {
    id: 'min_weight',
    name: 'min_weight',
    label: 'Minimum Weight',
    placeholder: 'Minimum weight of the selected asset in kg',
    value: '',
    type: 'number',
    size: 'small',
    isRequired: false,
    errorMsg: '',
    validators: [
      {
        validate: (value) => isNaN(value) || parseInt(value) <= 0,
        message: 'A positive number is required',
      },
      {
        validate: (value) => parseInt(value) > Number.MAX_SAFE_INTEGER,
        message: 'Value cannot exceed limit',
      },
    ],
  },
  max_height: {
    id: 'max_height',
    name: 'max_height',
    label: 'Maximum height',
    placeholder: 'Maximum height of selected asset in inches',
    value: '',
    type: 'number',
    size: 'small',
    isRequired: false,
    errorMsg: '',
    validators: [
      {
        validate: (value) => isNaN(value) || parseInt(value) <= 0,
        message: 'A positive number is required',
      },
      {
        validate: (value) => parseInt(value) > Number.MAX_SAFE_INTEGER,
        message: 'Value cannot exceed limit',
      },
    ],
  },
  min_height: {
    id: 'min_height',
    name: 'min_height',
    label: 'Minimum height',
    placeholder: 'Minimum height of selected asset in inches',
    value: '',
    type: 'number',
    size: 'small',
    isRequired: false,
    errorMsg: '',
    validators: [
      {
        validate: (value) => isNaN(value) || parseInt(value) <= 0,
        message: 'A positive number is required',
      },
      {
        validate: (value) => parseInt(value) > Number.MAX_SAFE_INTEGER,
        message: 'Value cannot exceed limit',
      },
    ],
  },
  return_notes: {
    id: 'return_notes',
    name: 'return_notes',
    label: 'Return notes',
    placeholder: 'Add return notes in less than 500 characters',
    value: '',
    type: 'text',
    isRequired: false,
    size: 'small',
    errorMsg: '',
    validators: [
      {
        validate: (value) => value.trim().length >= 500,
        message: 'Return notes should be less than 500 characters',
      },
    ],
  },
  created_at: {
    id: 'created_at',
    errorMsg: '',
    validators: [],
  },
  created_by: {
    id: 'created_by',
    errorMsg: '',
    validators: [],
  },
  updated_at: {
    id: 'updated_at',
    errorMsg: '',
    validators: [],
  },
  updated_by: {
    id: 'updated_by',
    errorMsg: '',
    validators: [],
  },
  sharable_groups: {
    id: 'sharable_groups',
    errorMsg: '',
    validators: [],
  },
};
