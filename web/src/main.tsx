import React, { useEffect, useState } from 'react';
import { createRoot } from 'react-dom/client';
import { Cloud, CloudRain, CloudSun, Compass, Droplets, Gauge, MapPin, RefreshCw, Search, Sun, Wind } from 'lucide-react';
import { getOverview, searchCities } from './services/api';
import type { City, Overview } from './types/weather';
import './style.css';
const defaultCity: City = { name: 'Prague', country: 'Czechia', latitude: 50.0755, longitude: 14.4378, timezone: 'Europe/Prague' };
function weatherIcon(code: number) { if (code >= 51) return CloudRain; if (code >= 1) return CloudSun; return Sun; }
function App() {
 const [city,setCity]=useState<City>(defaultCity),[data,setData]=useState<Overview|null>(null),[query,setQuery]=useState(''),[results,setResults]=useState<City[]>([]),[loading,setLoading]=useState(true),[error,setError]=useState(''),[unit,setUnit]=useState<'C'|'F'>('C');
 const load=async(c=city)=>{setLoading(true);setError('');try{setData(await getOverview(c));}catch(e){setError(e instanceof Error?e.message:'Weather request failed');}finally{setLoading(false);}};
 useEffect(()=>{void load(city);},[city]);
 useEffect(()=>{if(query.trim().length<2){setResults([]);return;}let active=true;const timer=window.setTimeout(()=>{searchCities(query.trim()).then(r=>{if(active)setResults(r);}).catch(()=>{if(active)setResults([]);});},300);return()=>{active=false;clearTimeout(timer);};},[query]);
 const temp=(c:number)=>Math.round(unit==='C'?c:c*9/5+32);const Icon=data?weatherIcon(data.current.weatherCode):CloudSun;
 const pick=(c:City)=>{setCity(c);setQuery('');setResults([]);};
 const widgets=data?[{label:'Feels like',value:`${temp(data.current.apparentTemperature)}°`,icon:CloudSun,note:'Apparent temperature'},{label:'Humidity',value:`${data.current.relativeHumidity}%`,icon:Droplets,note:'Relative humidity'},{label:'Wind',value:`${Math.round(data.current.windSpeed)} km/h`,icon:Wind,note:`Gusts ${Math.round(data.current.windGusts)} km/h`},{label:'Pressure',value:`${Math.round(data.current.pressure)} hPa`,icon:Gauge,note:'Sea level'}]:[];
 return <main className="app-shell"><header className="topbar"><a className="brand" href="#"><span className="brand-mark"><CloudSun size={21}/></span><span>WEATHER<span className="brand-light">/DASH</span></span></a><div className="search-wrap"><Search size={17}/><input aria-label="Search city" placeholder="Search city..." value={query} onChange={e=>setQuery(e.target.value)}/>{results.length>0&&<div className="search-results">{results.map((c,i)=><button key={`${c.latitude}-${i}`} onClick={()=>pick(c)}><MapPin size={15}/><span>{c.name}<small>{[c.region,c.country].filter(Boolean).join(', ')}</small></span></button>)}</div>}</div><div className="top-actions"><div className="unit-toggle"><button className={unit==='C'?'selected':''} onClick={()=>setUnit('C')}>°C</button><button className={unit==='F'?'selected':''} onClick={()=>setUnit('F')}>°F</button></div><button className="icon-button" aria-label="Refresh weather" onClick={()=>void load()}><RefreshCw size={17}/></button></div></header>
 <section className="welcome-row"><div><p className="eyebrow">WEATHER OVERVIEW</p><h1>Dashboard</h1></div><span className="live-label"><i/> LIVE CONDITIONS</span></section>
 {error&&<div className="error-banner">{error}<button onClick={()=>void load()}>Retry</button></div>}
 {loading&&!data?<section className="loading-card">Loading weather for {city.name}…</section>:data&&<>
 <section className="hero-card"><div className="hero-top"><div><div className="location"><MapPin size={16}/>{city.name}<span>{city.country}</span></div><p className="condition">{data.current.weatherCode===0?'Clear sky':data.current.weatherCode<4?'Partly cloudy':'Cloudy conditions'}</p></div><Icon className="hero-weather-icon" size={62}/></div><div className="hero-bottom"><div className="temperature">{temp(data.current.temperature)}<sup>°{unit}</sup></div><div className="hero-meta"><span>FEELS LIKE <b>{temp(data.current.apparentTemperature)}°</b></span><span>HUMIDITY <b>{data.current.relativeHumidity}%</b></span></div><div className="sun-times"><div><Sun size={16}/><span>Sunrise</span><b>{data.daily[0]?.sunrise.slice(-5)??'—'}</b></div><div><Compass size={16}/><span>Sunset</span><b>{data.daily[0]?.sunset.slice(-5)??'—'}</b></div></div></div></section>
 <section className="section-heading"><div><p className="eyebrow">RIGHT NOW</p><h2>Conditions</h2></div><span>Updated {new Date(data.updated).toLocaleTimeString([], {hour:'2-digit',minute:'2-digit'})}</span></section><section className="widget-grid">{widgets.map(w=><article className="widget" key={w.label}><div className="widget-heading"><span>{w.label}</span><w.icon size={17}/></div><strong>{w.value}</strong><small>{w.note}</small></article>)}</section>
 <section className="section-heading"><div><p className="eyebrow">NEXT HOURS</p><h2>Hourly forecast</h2></div></section><section className="hourly-strip">{data.hourly.slice(0,12).map((h,i)=>{const I=weatherIcon(h.weatherCode);return <article className={`hour-item ${i===0?'hour-now':''}`} key={h.time}><span>{i===0?'Now':new Date(h.time).toLocaleTimeString([], {hour:'numeric'})}</span><I size={19}/><b>{temp(h.temperature)}°</b><small><Droplets size={12}/>{h.precipitationProbability}%</small></article>})}</section>
 </>}
 <footer>Weather data by <a href="https://open-meteo.com/" target="_blank" rel="noreferrer">Open-Meteo</a><span>{data?.timezone??city.timezone}</span></footer>
 </main>
}
createRoot(document.getElementById('root')!).render(<React.StrictMode><App/></React.StrictMode>);
