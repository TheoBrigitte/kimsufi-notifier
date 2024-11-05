'use client'

import ServerLine from './components/server';
import useSWR from 'swr';

const getServers = () => fetch('http://127.0.0.1:8080/list').then(res => res.json());

export default function Home() {
  const { data, error, isLoading } = useSWR('/list', getServers);

  let content;

  if (error) {
    content =
    <div className="flex flex-col justify-center text-center">
      <div className="flex flex-row justify-center">
        <div className="w-1/2 text-right px-1">Status :</div>
        <div className="w-1/2 text-left px-1">Failed to load server list</div>
      </div>
      <div className="text-orange-700 font-mono">{error.toString()}</div>
    </div>;
  } else if (isLoading) {
    content =
    <div className="flex flex-col justify-center text-center">
      <div className="flex flex-row justify-center">
        <div className="w-1/2 text-right px-1">Status :</div>
        <div className="w-1/2 text-left px-1">Loading ...</div>
      </div>
    </div>;
  } else {
    const serversByCategory = Object.groupBy(data, ({ category }) => category);

    let categoryOrder = { "Kimsufi": {}, "So you Start": {}, "Rise": {}, "": {} }
    const ordered = Object.assign(categoryOrder, serversByCategory);
    console.log(ordered);

    content = <table className="text-nowrap">
        <thead>
          <tr>
            <th className="p-4">Plan Code</th>
            <th className="p-4">Category</th>
            <th className="p-4">Name</th>
            <th className="p-4">Price</th>
            <th className="p-4">Status</th>
            <th className="p-4">Datacenters</th>
          </tr>
        </thead>
        <tbody>
          {Object.entries(ordered).map(([category, servers]) => (
            <>
            <ServerLine key={category} servers={servers} />
            <tr key={category + " separator0"}><td className="p-2" colSpan={6}></td></tr>
            <tr key={category + " separator1"}><td className="p-2" colSpan={6}></td></tr>
            </>
          ))}
        </tbody>
      </table>
  }


  return (
    <div className="flex flex-row justify-center">
    <div className="pt-10 pb-20">
      <h1 className="text-center text-xl font-bold p-10">OVH Eco server availability</h1>
      {content}
    </div>
    </div>
  );
}
