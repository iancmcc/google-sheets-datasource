import { DataQuery, DataSourceJsonData, SelectableValue } from '@grafana/data';

//-------------------------------------------------------------------------------
// General google cloud auth types
// same as stackdriver etc
//-------------------------------------------------------------------------------

export interface JWT {
  private_key: any;
  token_uri: any;
  client_email: any;
  project_id: any;
}

export enum GoogleAuthType {
  JWT = 'jwt',
  GCE = 'gce',
  NONE = 'none',
}

export const googleAuthTypes = [
  { label: 'None (public)', value: GoogleAuthType.NONE },
  { label: 'Google JWT File', value: GoogleAuthType.JWT },
  { label: 'GCE Default Service Account', value: GoogleAuthType.GCE },
];

export interface GoogleCloudOptions extends DataSourceJsonData {
  authenticationType: GoogleAuthType;
}

//-------------------------------------------------------------------------------
// The Sheets specicif types
//-------------------------------------------------------------------------------
export interface GoogleSheetRangeInfo {
  spreadsheet: SelectableValue<string>;
  range: string;
}

export interface SheetsQuery extends DataQuery, GoogleSheetRangeInfo {
  cacheDurationSeconds: number;
  queryType: string;
}

export interface SheetsSourceOptions extends GoogleCloudOptions {
  authType: GoogleAuthType;
  apiKey: string;
  jwtFile: string;
}

export interface GoogleSheetsSecureJsonData {
  apiKey: string;
}
