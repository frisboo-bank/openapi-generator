import { listEntities } from "../../../../grpc/client/entity/v1/entity_service-EntityService_connectquery"
import { useQuery } from '@connectrpc/connect-query';

function EntitiesList() {
  const { data } = useQuery(listEntities, {

  })

  console.error(data);

  return <div>Hello</div>
}

export default EntitiesList;
