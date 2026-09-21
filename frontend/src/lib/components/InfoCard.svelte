<script lang="ts">
  import { ChevronRight } from "lucide-svelte";
  import { cn } from "$lib/utils";
  import { i18n } from "$lib/i18n.svelte";

  let { category, title, summary, date, source, icon: Icon = null, accent = "", href = "#", onclick } = $props<{
    category: string,
    title: string,
    summary: string,
    date: string,
    source: string,
    icon?: any,
    accent?: string,
    href?: string,
    onclick?: (e: MouseEvent) => void
  }>();

  const CATEGORY_STYLES: Record<string, string> = {
    volcano: "bg-warning/10 text-warning",
    lahar: "bg-[#92400E]/10 text-[#92400E]",
    evac: "bg-primary/10 text-primary",
    weather: "bg-secondary text-white",
    warning: "bg-destructive/10 text-destructive",
    default: "bg-muted text-muted-foreground",
  };

  let catStyle = $derived(CATEGORY_STYLES[category] || CATEGORY_STYLES.default);
</script>

<a {href} onclick={(e) => { if(onclick) { e.preventDefault(); onclick(e); } }} class="group flex flex-col p-5 rounded-xl bg-card border border-border hover:shadow-[0_4px_12px_-2px_rgba(15,23,42,0.06)] transition-shadow text-left">
  <div class="flex items-center justify-between">
    {#if Icon}
      <span class={cn("grid place-items-center w-10 h-10 rounded-xl", accent || "bg-primary/10 text-primary")}>
        <Icon class="w-5 h-5" strokeWidth={2} />
      </span>
    {/if}
    <span class={cn("text-[0.75rem] font-semibold px-2.5 py-1 rounded-full", catStyle)}>{i18n.t(`info.filter.${category}` as any) || category}</span>
  </div>
  <h3 class="mt-4 font-display font-semibold text-[1.0625rem] text-foreground leading-snug">{title}</h3>
  <p class="mt-2 text-[0.9375rem] text-muted-foreground leading-relaxed line-clamp-2">{summary}</p>
  <div class="mt-4 flex flex-col gap-1.5 text-[0.8125rem] text-muted-foreground">
    <span class="font-medium text-foreground">{date}</span>
    <span>{source}</span>
  </div>
  <span class="mt-4 inline-flex items-center gap-1 text-[0.8125rem] font-semibold text-primary group-hover:gap-1.5 transition-all self-start">
    {i18n.t('info.card.readmore')} <ChevronRight class="w-3.5 h-3.5" />
  </span>
</a>
