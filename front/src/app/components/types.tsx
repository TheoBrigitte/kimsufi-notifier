import { IconDefinition } from '@fortawesome/free-regular-svg-icons';

type Server = {
  bandwidth: string;
  category: string;
  cpu: string;
  currencyCode: string;
  datacenters: string[];
  memory: string;
  name: string;
  planCode: string;
  price: number;
  status: string;
  storage: string;
};

type Status = {
  color: string;
  icon: IconDefinition;
};

export type { Server, Status };
